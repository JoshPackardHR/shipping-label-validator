package shipping

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"strings"

	"github.com/JoshuaPackardHR/shipping-label-validator/gpt"
	"github.com/JoshuaPackardHR/shipping-label-validator/internal/shipping/models"
	"github.com/JoshuaPackardHR/shipping-label-validator/ups"
)

//go:embed prompt.txt
var prompt string

type manager struct {
	upsClient ups.Client
	gpt       gpt.GPT
}

func NewManager(
	upsClient ups.Client,
	gpt gpt.GPT,
) models.Manager {
	return &manager{
		upsClient: upsClient,
		gpt:       gpt,
	}
}

type promptResponse struct {
	ups.Address
	TrackingNumber string `json:"trackingNumber"`
	Error          string `json:"error"`
}

// TODO: Store validation results in the database
func (m *manager) Validate(ctx context.Context, trackingNumber string, img image.Image) (*models.ValidationResult, error) {
	imageBytes := new(bytes.Buffer)
	if err := jpeg.Encode(imageBytes, img, nil); err != nil {
		return nil, err
	}

	// Call LLM to read the address and tracking number from the image
	result, err := m.gpt.Prompt(ctx, prompt, imageBytes.Bytes())
	if err != nil {
		return nil, err
	}
	promptResp := promptResponse{}
	if err := json.Unmarshal([]byte(result.Content), &promptResp); err != nil {
		return nil, err
	}

	// Call UPS API to get the address for the tracking number
	if trackingNumber == "" {
		trackingNumber = promptResp.TrackingNumber
	}
	trackingDetails, err := m.upsClient.GetTrackingDetails(trackingNumber)
	if err != nil {
		return nil, err
	}
	expectedAddress := trackingDetails.GetPackageAddress(ups.PackageAddressTypeDestination)
	if expectedAddress == nil {
		return nil, errors.New("no address found for the tracking number")
	}

	// Compare the address from the image with the address from the UPS API
	return &models.ValidationResult{
		ScannedAddress:         promptResp.Address,
		ExpectedPackageAddress: *expectedAddress,
		Valid:                  compareAddresses(promptResp.Address, expectedAddress.Address),
	}, nil
}

func compareAddresses(address1, address2 ups.Address) bool {
	return strings.EqualFold(address1.AddressLine1, address2.AddressLine1) &&
		strings.EqualFold(address1.AddressLine2, address2.AddressLine2) &&
		strings.EqualFold(address1.City, address2.City) &&
		strings.EqualFold(address1.StateProvince, address2.StateProvince) &&
		strings.EqualFold(address1.PostalCode, address2.PostalCode)
}
