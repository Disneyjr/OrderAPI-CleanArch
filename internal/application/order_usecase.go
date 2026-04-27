package application

import (
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/interfaces"

	"github.com/Disneyjr/OrderAPI-CleanArch/internal/domain"
)

type OrderInputDTO struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Tax         float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         uint64  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type OrderUseCase struct {
	orderRepository interfaces.OrderRepository
}

func NewOrderUseCase(orderRepository interfaces.OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		orderRepository: orderRepository,
	}
}

func (u *OrderUseCase) GetAll() ([]*domain.Order, error) {
	res, err := u.orderRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (u *OrderUseCase) Get(i uint64) (*domain.Order, error) {
	res, err := u.orderRepository.Get(i)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *OrderUseCase) Create(input OrderInputDTO) (OrderOutputDTO, error) {
	order := domain.Order{
		Description: input.Description,
		Price:       input.Price,
		Tax:         input.Tax,
	}
	order.CalculateFinalPrice()
	if err := order.Validate(); err != nil {
		return OrderOutputDTO{}, err
	}
	res, err := u.orderRepository.Post(&order)
	if err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         res.ID,
		Price:      res.Price,
		Tax:        res.Tax,
		FinalPrice: res.FinalPrice,
	}

	return dto, nil
}
