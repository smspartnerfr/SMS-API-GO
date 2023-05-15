package smspartner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type NumberVerificationRequest struct {
	APIKey       string `json:"apiKey,omitempty"`
	PhoneNumbers string `json:"phoneNumbers,omitempty"`
	NotifyURL    string `json:"notifyUrl,omitempty"`
}

type NumberVerificationResponse struct {
	Success    bool    `json:"success,omitempty"`
	Code       int     `json:"code,omitempty"`
	CampaignID string  `json:"campaign_id,omitempty"`
	Number     int     `json:"number,omitempty"`
	Cost       float64 `json:"cost,omitempty"`
	Currency   string  `json:"currency,omitempty"`
}

// VerifyNumber checks that a phone number actually exists.
func (c *Client) VerifyNumber(reqPayload *NumberVerificationRequest) (*NumberVerificationResponse, error) {
	reqPayload.APIKey = c.apiKey

	blob, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s/hlr/notify", c.basePath)

	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	blob, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	nvr := new(NumberVerificationResponse)
	if err := json.Unmarshal(blob, &nvr); err != nil {
		return nil, err
	}
	return nvr, nil
}

type NumberFormat struct {
	E164          string `json:"e164,omitempty"`
	International string `json:"international,omitempty"`
	National      string `json:"national,omitempty"`
	RFC3966       string `json:"rfc3966,omitempty"`
}

type Lookup struct {
	Request     string        `json:"request,omitempty"`
	Success     bool          `json:"success,omitempty"`
	CountryCode string        `json:"countryCode,omitempty"`
	PrefixCode  int           `json:"prefixCode,omitempty"`
	PhoneNumber string        `json:"phoneNumber,omitempty"`
	Type        string        `json:"type,omitempty"`
	Network     string        `json:"network,omitempty"`
	Format      *NumberFormat `json:"format,omitempty"`
}

type LookupResponse struct {
	Success bool      `json:"success,omitempty"`
	Code    int       `json:"code,omitempty"`
	Lookup  []*Lookup `json:"lookup,omitempty"`
}

// VerifyNumberFormat checks the format of a phone number
func (c *Client) VerifyNumberFormat(phoneNumbers ...string) (*LookupResponse, error) {
	if len(phoneNumbers) == 0 {
		return nil, errors.New("At least one phoneNumber is required")
	}
	p := strings.Join(phoneNumbers, ",")

	payload := new(NumberVerificationRequest)
	payload.APIKey = c.apiKey
	payload.PhoneNumbers = p

	blob, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s/lookup", c.basePath)

	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	blob, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	lr := new(LookupResponse)
	if err := json.Unmarshal(blob, lr); err != nil {
		return nil, err
	}
	return lr, nil
}
