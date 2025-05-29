package models

import (
	"context"
	"image"

	"github.com/JoshuaPackardHR/shipping-label-validator/ups"
)

type Manager interface {
	Validate(ctx context.Context, trackingNumber string, image image.Image) (*ValidationResult, error)
}

type ValidationResult struct {
	ScannedAddress         ups.Address        `json:"scannedAddress"`
	ExpectedPackageAddress ups.PackageAddress `json:"expectedAddress"`
	Valid                  bool               `json:"valid"`
} // @name ValidationResult
