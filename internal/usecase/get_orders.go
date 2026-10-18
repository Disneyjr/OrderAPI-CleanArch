package usecase

import "github.com/Disneyjr/OrderAPI-CleanArch/internal/entity"

type GetOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewGetOrdersUseCase(OrderRepository entity.OrderRepositoryInterface) *GetOrdersUseCase {
	return &GetOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}
type ListOrdersUseCase = GetOrdersUseCase

func NewListOrdersUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return NewGetOrdersUseCase(OrderRepository)
}

func (c *GetOrdersUseCase) Execute() ([]OrderOutputDTO, error) {
	ordersDTO := []OrderOutputDTO{}
	orders, err := c.OrderRepository.GetAll()
	if err != nil {
		return ordersDTO, err
	}
	for _, order := range orders {
		dto := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		}
		ordersDTO = append(ordersDTO, dto)
	}
	return ordersDTO, nil
}

