package handler

import (
	"context"
	"inventory-service/internal/model"
	"inventory-service/internal/usecase"
	pb "inventory-service/pb/inventory"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	pb.UnimplementedInventoryServiceServer
	Usecase *usecase.ProductUsecase
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	product := model.Product{
		Name:     req.Name,
		Category: req.Category,
		Stock:    int(req.Stock),
		Price:    req.Price,
	}

	err := h.Usecase.Create(&product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &pb.Product{
		Id:       int64(product.ID),
		Name:     product.Name,
		Category: product.Category,
		Stock:    int32(product.Stock),
		Price:    product.Price,
	}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.ProductID) (*pb.Product, error) {
	product, err := h.Usecase.GetByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	return &pb.Product{
		Id:       int64(product.ID),
		Name:     product.Name,
		Category: product.Category,
		Stock:    int32(product.Stock),
		Price:    product.Price,
	}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	product := model.Product{
		ID:       int(req.Id),
		Name:     req.Name,
		Category: req.Category,
		Stock:    int(req.Stock),
		Price:    req.Price,
	}

	err := h.Usecase.Update(product.ID, &product) // ✅ fixed
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &pb.Product{
		Id:       int64(product.ID),
		Name:     product.Name,
		Category: product.Category,
		Stock:    int32(product.Stock),
		Price:    product.Price,
	}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.ProductID) (*pb.Empty, error) {
	err := h.Usecase.Delete(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}
	return &pb.Empty{}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, _ *pb.Empty) (*pb.ProductList, error) {
	products, err := h.Usecase.List("", 100, 0) // or implement filtering later
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	var protoProducts []*pb.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.Product{
			Id:       int64(p.ID),
			Name:     p.Name,
			Category: p.Category,
			Stock:    int32(p.Stock),
			Price:    p.Price,
		})
	}

	return &pb.ProductList{Products: protoProducts}, nil
}
