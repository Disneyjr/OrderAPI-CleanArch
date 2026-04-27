package postgres

import (
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/domain"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}
func (r *OrderRepository) GetAll() ([]*domain.Order, error) {
	query := r.db.Find(&[]*domain.Order{})
	if query.Error != nil {
		return nil, query.Error
	}
	orders := []*domain.Order{}
	if err := query.Scan(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) Get(i uint64) (*domain.Order, error) {
	query := r.db.Find(&domain.Order{ID: i})
	if query.Error != nil {
		return nil, query.Error
	}
	order := &domain.Order{}
	if err := query.Scan(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) Post(order *domain.Order) (*domain.Order, error) {
	query := r.db.Create(order)
	if query.Error != nil {
		return nil, query.Error
	}
	return order, nil
}
