package models

import (
	"context"
	"image"

	"github.com/happyreturns/shipping-label-checker/ups"
)

type Manager interface {
	Validate(ctx context.Context, trackingNumber string, image image.Image) (*ValidationResult, error)
}

type ValidationResult struct {
	ScannedAddress  ups.Address `json:"scannedAddress"`
	ExpectedAddress ups.Address `json:"expectedAddress"`
	Valid           bool        `json:"valid"`
} // @name ValidationResult
