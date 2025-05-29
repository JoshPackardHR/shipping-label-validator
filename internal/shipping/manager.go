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

	"github.com/happyreturns/shipping-label-checker/gpt"
	"github.com/happyreturns/shipping-label-checker/internal/shipping/models"
	"github.com/happyreturns/shipping-label-checker/ups"
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
	ShipTo         string `json:"shipTo"`
	TrackingNumber string `json:"trackingNumber"`
	Error          string `json:"error"`
}

// TODO: Store results in the database
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
	// TODO: Parse scanned address better
	scannedAddress := ups.Address{AddressLine1: promptResp.ShipTo}

	// Call UPS API to get the address for the tracking number
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
		ScannedAddress:  scannedAddress,
		ExpectedAddress: *expectedAddress,
		Valid:           compareAddresses(scannedAddress, *expectedAddress),
	}, nil
}

func compareAddresses(address1, address2 ups.Address) bool {
	// TODO: Compare addresses better
	return strings.Contains(strings.ToLower(address1.AddressLine1), strings.ToLower(address2.AddressLine1))
}
