package interfaces

import (
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/domain"
)

type OrderRepository interface {
	GetAll() ([]*domain.Order, error)
	Get(uint64) (*domain.Order, error)
	Post(*domain.Order) (*domain.Order, error)
}
