package database

import (
	"database/sql"

	"github.com/Disneyjr/OrderAPI-CleanArch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetAll() ([]entity.Order, error) {
	var orders []entity.Order
	rows, err := r.Db.Query("SELECT * FROM ORDERS")
	if err != nil {
		return orders, err
	}

	if err := rows.Scan(orders); err != nil {
		return orders, err
	}
	return orders, nil
}
