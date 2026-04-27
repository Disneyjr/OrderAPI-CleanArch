package domain

import "errors"

type Order struct {
	ID          uint64  `json:"id" gorm:"primaryKey"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Tax         float64 `json:"tax"`
	FinalPrice  float64 `json:"final_price"`
}

func (o *Order) Validate() error {
	if o.Price <= 0 {
		return errors.New("invalid price")
	}
	if o.Tax <= 0 {
		return errors.New("invalid tax")
	}
	return nil
}
func (o *Order) CalculateFinalPrice() {
	o.FinalPrice = o.Price + o.Tax
}
