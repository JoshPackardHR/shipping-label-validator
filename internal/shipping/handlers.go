package shipping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/happyreturns/shipping-label-checker/internal/shipping/models"
	"github.com/inpersondonations/helpers/api"
)

type handler struct {
	manager models.Manager
}

func NewHandler(manager models.Manager) *handler {
	return &handler{manager: manager}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/label/check", h.login)
}

type ShippingRequest struct {
	Image string `json:"image"`
} // @name ShippingRequest

type ShippingResponse struct {
	Valid bool `json:"valid"`
} // @name ShippingResponse

type ShippingError struct {
	Error string `json:"error"`
} // @name ShippingError

// login godoc
//
//	@Summary		Check a shipping label
//	@Description	check a shipping label
//	@Tags			shipping
//	@Accept			json
//	@Produce		json
//	@Param			requestBody	body		ShippingRequest	true	"Shipping Label Request"
//	@Success		200			{object}	ShippingResponse
//	@Failure		400,500		{object}	ShippingError
//	@Router			/shipping/label/check [post]
func (h *handler) login(c *gin.Context) {
	shippingRequest := ShippingRequest{}
	if err := c.BindJSON(&shippingRequest); err != nil {
		c.JSON(http.StatusBadRequest, ShippingError{Error: err.Error()})
		return
	}

	valid, err := h.manager.CheckLabel(c, shippingRequest.Image)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, ShippingResponse{Valid: valid})
}
