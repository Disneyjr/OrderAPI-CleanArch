package service

import (
	"context"

	"github.com/Disneyjr/OrderAPI-CleanArch/internal/infra/grpc/pb"
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	GetOrderUseCase    usecase.GetOrdersUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, getOrderUseCase usecase.GetOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		GetOrderUseCase:    getOrderUseCase,
	}
}
func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

/*
func (s *OrderService) GetOrder(ctx context.Context) (*pb.GetOrderResponse, error) {
	output, err := s.GetOrderUseCase.Execute()
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}
*/
