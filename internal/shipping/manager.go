package shipping

import (
	"context"

	"github.com/happyreturns/shipping-label-checker/internal/shipping/models"
)

type manager struct {
}

func NewManager() models.Manager {
	return &manager{}
}
func (m *manager) CheckLabel(ctx context.Context, image string) (bool, error) {
	// TODO: Call OpenAI API to read the address and tracking number from the image
	// TODO: UPS API integration to check the address for the tracking number
	// TODO: Compare the address from the image with the address from the UPS API
	return true, nil
}
