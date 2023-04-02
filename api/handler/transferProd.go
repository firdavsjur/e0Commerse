package handler

import (
	"app/api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Transfer Product godoc
// @ID tranfer_product
// @Router /transferProduct [PUT]
// @Summary transfer transferProduct
// @Description transfer product
// @Tags TransferProduct
// @Accept json
// @Produce json
// @Param transferProduct body models.TransferProductRequest true "transferProudctRequest"
// @Success 201 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) TransferProduct(c *gin.Context) {
	var transferPr models.TransferProductRequest

	err := c.ShouldBindJSON(&transferPr) // parse req body to given type struct
	if err != nil {
		h.handlerResponse(c, "transfer product", http.StatusBadRequest, err.Error())
		return
	}

	err = h.storages.Stock().TransferProduct(context.Background(), &transferPr)
	if err != nil {
		h.handlerResponse(c, "transer product exec", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "transferProduct", http.StatusCreated, "transfered")
}
