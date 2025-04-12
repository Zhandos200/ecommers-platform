package handler

import (
	"context"
	"order-service/internal/model"
	"order-service/internal/usecase"

	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "order-service/pb/order"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	usecase *usecase.OrderUsecase
}

// handler/order_handler.go
func NewOrderHandler(uc *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: uc}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	order := &model.Order{
		UserID: int(req.UserId), // fixed: proto = int64 → model = int
		Status: "pending",
	}

	for _, item := range req.Items {
		order.Items = append(order.Items, model.OrderItem{
			ProductID: int(item.ProductId),
			Quantity:  int(item.Quantity),
		})
	}

	if err := h.usecase.Create(order); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return convertToOrderResponse(order), nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.OrderID) (*pb.OrderResponse, error) {
	order, err := h.usecase.GetByID(int(req.Id)) // order is *model.Order
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}
	return convertToOrderResponse(order), nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.StatusUpdate) (*pb.OrderResponse, error) {
	if err := h.usecase.UpdateStatus(int(req.Id), req.Status); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status")
	}
	order, err := h.usecase.GetByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}
	return convertToOrderResponse(order), nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.UserOrdersRequest) (*pb.OrderList, error) {
	var (
		orders []model.Order
		err    error
	)

	if req.UserId != 0 {
		orders, err = h.usecase.ListByUser(int(req.UserId))
	} else {
		orders, err = h.usecase.ListAll()
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	var result pb.OrderList
	for i := range orders {
		result.Orders = append(result.Orders, convertToOrderResponse(&orders[i]))
	}
	return &result, nil

}

func convertToOrderResponse(order *model.Order) *pb.OrderResponse {
	var items []*pb.OrderItem
	for _, item := range order.Items {
		items = append(items, &pb.OrderItem{
			ProductId: int64(item.ProductID),
			Quantity:  int32(item.Quantity),
		})
	}
	return &pb.OrderResponse{
		Id:        int64(order.ID),
		UserId:    int64(order.UserID),
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
		Items:     items,
	}

}
