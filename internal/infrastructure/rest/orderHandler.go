package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	order_usecase "github.com/Disneyjr/OrderAPI-CleanArch/internal/application"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	order_usecase *order_usecase.OrderUseCase
}

func NewOrderHandler(order_usecase *order_usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{order_usecase: order_usecase}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	resp, err := h.order_usecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	num, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	resp, err := h.order_usecase.Get(num)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func (h *OrderHandler) Create(c *gin.Context) {
	r := c.Request
	var dto order_usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		c.Error(err)
		return
	}
	output, err := h.order_usecase.Create(dto)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusCreated, output)

}
