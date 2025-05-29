package shipping

import (
	"encoding/base64"
	"image/jpeg"
	"net/http"
	"strings"

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
	router.POST("/label/validate", h.validate)
}

type ValidationRequest struct {
	TrackingNumber string `json:"trackingNumber"`
	Image          string `json:"image"`
} // @name ValidationRequest

type ValidationResponse struct {
	Result models.ValidationResult `json:"result"`
} // @name ValidationResponse

type ValidationError struct {
	Error string `json:"error"`
} // @name ValidationError

// login godoc
//
//	@Summary		Check a shipping label
//	@Description	check a shipping label
//	@Tags			shipping
//	@Accept			json
//	@Produce		json
//	@Param			requestBody	body		ValidationRequest	true	"Validation Request"
//	@Success		200			{object}	ValidationResponse
//	@Failure		400,500		{object}	ValidationError
//	@Router			/shipping/label/validate [post]
func (h *handler) validate(c *gin.Context) {
	request := ValidationRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ValidationError{Error: err.Error()})
		return
	}

	// parse base64 image to native image
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(request.Image))
	image, err := jpeg.Decode(reader)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	result, err := h.manager.Validate(c, request.TrackingNumber, image)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, ValidationResponse{Result: *result})
}
