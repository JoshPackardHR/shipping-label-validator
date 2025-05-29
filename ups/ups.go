package ups

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	timeoutDuration       = 30 * time.Second
	tlsHandshakeTimeout   = 10 * time.Second
	idleConnTimeout       = 10 * time.Second
	responseHeaderTimeout = 10 * time.Second
	expectContinueTimeout = 10 * time.Second
	tokenUrl              = "https://onlinetools.ups.com/security/v1/oauth/token"
	trackingUrl           = "https://onlinetools.ups.com/api/track/v1/details"
)

type Client interface {
	GetTrackingDetails(trackingNumber string) (*TrackingDetails, error)
}

type client struct {
	accessToken string
}

func NewClient(clientId string, clientSecret string) (Client, error) {
	token, err := getAccessToken(nil, clientId, clientSecret, nil, nil)
	if err != nil {
		return nil, err
	}

	return &client{
		accessToken: token.AccessToken,
	}, nil
}

type PackageAddressType string // @name PackageAddressType

const PackageAddressTypeOrigin = "ORIGIN"
const PackageAddressTypeDestination = "DESTINATION"

type Address struct {
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
	City          string `json:"city"`
	StateProvince string `json:"stateProvince"`
	PostalCode    string `json:"postalCode"`
	CountryCode   string `json:"countryCode"`
	Country       string `json:"country"`
} // @name Address

type PackageAddress struct {
	Type          PackageAddressType `json:"type"`
	Name          string             `json:"name"`
	AttentionName string             `json:"attentionName"`
	Address       Address            `json:"address"`
} // @name PackageAddress

type TrackingDetails struct {
	TrackResponse struct {
		Shipment []struct {
			InquiryNumber string `json:"inquiryNumber"`
			ShipmentType  string `json:"shipmentType"`
			ShipperNumber string `json:"shipperNumber"`
			PickupDate    string `json:"pickupDate"`
			Package       []struct {
				TrackingNumber string           `json:"trackingNumber"`
				PackageAddress []PackageAddress `json:"packageAddress"`
			} `json:"package"`
			UserRelation []string `json:"userRelation"`
		} `json:"shipment"`
	} `json:"trackResponse"`
}

func (t TrackingDetails) GetPackageAddress(addressTypeType PackageAddressType) *PackageAddress {
	for _, ship := range t.TrackResponse.Shipment {
		for _, pkg := range ship.Package {
			for _, addr := range pkg.PackageAddress {
				if addr.Type == addressTypeType {
					return &addr
				}
			}
		}
	}

	return nil
}

func (c *client) GetTrackingDetails(trackingNumber string) (*TrackingDetails, error) {
	var hClient *http.Client = setHttpClientTimeouts(nil)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", trackingUrl, trackingNumber), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("transId", fmt.Sprintf("%d", time.Now().Unix()))
	req.Header.Set("transactionSrc", "testing")

	res, err := hClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, err
		}
		return nil, err
	}

	defer res.Body.Close()

	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	respString := string(response)
	var data TrackingDetails
	err = json.Unmarshal([]byte(respString), &data)
	if err != nil {
		return nil, err
	}

	if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
		return nil, errors.New(respString)
	}

	return &data, nil
}

func setHttpClientTimeouts(httpClient *http.Client) *http.Client {
	if httpClient == nil {
		return &http.Client{
			Timeout: timeoutDuration,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   tlsHandshakeTimeout,
				IdleConnTimeout:       idleConnTimeout,
				ResponseHeaderTimeout: responseHeaderTimeout,
				ExpectContinueTimeout: expectContinueTimeout,
			},
		}
	}

	httpClient.Timeout = timeoutDuration

	if transport, ok := httpClient.Transport.(*http.Transport); ok {
		transport.TLSHandshakeTimeout = tlsHandshakeTimeout
		transport.IdleConnTimeout = idleConnTimeout
		transport.ResponseHeaderTimeout = responseHeaderTimeout
		transport.ExpectContinueTimeout = expectContinueTimeout
	} else {
		httpClient.Transport = &http.Transport{
			TLSHandshakeTimeout:   tlsHandshakeTimeout,
			IdleConnTimeout:       idleConnTimeout,
			ResponseHeaderTimeout: responseHeaderTimeout,
			ExpectContinueTimeout: expectContinueTimeout,
		}
	}
	return httpClient
}

type TokenInfo struct {
	IssuedAt    string `json:"issued_at"`
	TokenType   string `json:"token_type"`
	ClientId    string `json:"client_id"`
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	Status      string `json:"status"`
}

func getAccessToken(httpClient *http.Client, clientId string, clientSecret string, headers map[string]string, customClaims map[string]string) (*TokenInfo, error) {
	var hClient *http.Client = setHttpClientTimeouts(httpClient)

	body := url.Values{}
	body.Set("grant_type", "client_credentials")
	body.Set("scope", "public")

	for keys := range customClaims {
		claim := "{\"" + keys + "\":\"" + customClaims[keys] + "\"}"
		body.Add("custom_claims", claim)
	}
	encodedData := body.Encode()

	req, err := http.NewRequest(http.MethodPost, tokenUrl, strings.NewReader(encodedData))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(clientId, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for keys := range headers {
		req.Header.Set(keys, headers[keys])
	}

	res, err := hClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, err
		}
		return nil, err
	}

	defer res.Body.Close()

	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	respString := string(response)
	var data TokenInfo
	errr := json.Unmarshal([]byte(respString), &data)
	if errr != nil {
		return nil, err
	}

	if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
		return nil, errors.New(respString)
	}

	return &data, nil
}
