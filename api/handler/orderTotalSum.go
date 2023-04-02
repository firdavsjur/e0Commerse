package handler

import (
	"app/api/models"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get By ID OrderTotal godoc
// @ID order_total
// @Router /orderTotal/{order_id} [GET]
// @Summary Get OrderTotal
// @Description  OrderTotal
// @Tags OrderTotal
// @Accept json
// @Produce json
// @Param order_id path string true "order_id"
// @Param promo_code query string false "promo_code"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) GetOrderTotalSum(c *gin.Context) {

	orderId := c.Param("id")
	promo_code_name := c.Query("promo_code")

	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		h.handlerResponse(c, "storage.orderTotal.getByID", http.StatusBadRequest, "id incorrect")
		return
	}

	resp, message, err := h.storages.Order().TotalOrderSum(context.Background(), &models.TotalOrderSumRequest{OrderId: orderIdInt, PromoCodeName: promo_code_name})
	if err != nil {
		h.handlerResponse(c, message, http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, message, http.StatusCreated, resp)
}
