package handler

import (
	"context"
	"user-service/internal/model"
	"user-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "user-service/pb/user" // Adjust this path to match your actual proto module
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *pb.UserRequest) (*pb.RegisterUserResponse, error) {
	user := &model.User{Email: req.Email, Name: req.Name, Password: req.Password}
	if err := h.usecase.Register(user); err != nil {
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}
	return &pb.RegisterUserResponse{Message: "Verification email sent"}, nil
}

func (h *UserHandler) AuthenticateUser(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	user, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: "Invalid credentials",
		}, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	return &pb.AuthResponse{
		Success: true,
		Message: "Login successful",
		User: &pb.UserResponse{
			Id:    int64(user.ID),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (h *UserHandler) GetUserProfile(ctx context.Context, req *pb.UserID) (*pb.UserProfile, error) {
	user, err := h.usecase.GetProfile(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.UserProfile{
		Id:    int64(user.ID),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user := model.User{
		ID:       int(req.Id),
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password, // optional: ignore if not updating password
	}

	if err := h.usecase.UpdateUser(user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UserResponse{
		Id:    int64(user.ID),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.UserID) (*emptypb.Empty, error) {
	if err := h.usecase.DeleteUser(int(req.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return &emptypb.Empty{}, nil
}
func (h *UserHandler) VerifyUser(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	if err := h.usecase.Verify(req.Token); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "verify failed: %v", err)
	}
	return &pb.VerifyResponse{Success: true}, nil
}
