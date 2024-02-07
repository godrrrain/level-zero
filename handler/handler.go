package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"zero/model"
	"zero/storage"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	storage storage.Storage
	cache   map[string]*model.OrderDb
}

func NewHandler(storage storage.Storage, cache map[string]*model.OrderDb) *Handler {
	return &Handler{
		storage: storage,
		cache:   cache,
	}
}

func (h *Handler) GetOrderByUid(c *gin.Context) {

	orderUID := c.Param("uid")

	// orderDb, err := h.storage.GetOrderDb(context.TODO(), orderUID)
	// if err != nil {
	// 	println(err)
	// }

	orderdb, ok := h.cache[orderUID]
	if !ok {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: fmt.Sprintf("not found %s", orderUID),
		})
		return
	}

	c.JSON(http.StatusOK, OrderDBToOrder(orderdb))
}

func OrderDBToOrder(orderdb *model.OrderDb) model.Order {
	var order model.Order
	err := json.Unmarshal(orderdb.DataJson, &order)
	if err != nil {
		fmt.Println(err)
	}
	return order
}
