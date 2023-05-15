package smspartner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Gamme is the SMS range to specify when sending SMS
type Gamme int

// List of values that Gamme can take.
const (
	Premium Gamme = 1
	LowCost Gamme = 2
)

type SMSPayload struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Message     string `json:"message,omitempty"`
}

type SMS struct {
	APIKey                string `json:"apiKey,omitempty"`
	PhoneNumbers          string `json:"phoneNumbers,omitempty"`
	Message               string `json:"message,omitempty"`
	Gamme                 Gamme  `json:"gamme,omitempty"`
	Sender                string `json:"sender,omitempty"`
	ScheduledDeliveryDate string `json:"scheduledDeliveryDate,omitempty"`
	Time                  int    `json:"time,omitempty"`
	Minute                int    `json:"minute,omitempty"`
	// IsStopSms
	// Sandbox
}

type BulkSMS struct {
	APIKey                string        `json:"apiKey,omitempty"`
	SMSList               []*SMSPayload `json:"SMSList,omitempty"`
	Gamme                 Gamme         `json:"gamme,omitempty"`
	Sender                string        `json:"sender,omitempty"`
	ScheduledDeliveryDate string        `json:"scheduledDeliveryDate,omitempty"`
	Time                  int           `json:"time,omitempty"`
	Minute                int           `json:"minute,omitempty"`
	// IsStopSms
	// Sandbox
}

type SMSResponse struct {
	Success               bool    `json:"success"`
	Code                  int     `json:"code"`
	MessageID             int     `json:"message_id"`
	NumberOfSMS           int     `json:"nb_sms"`
	Cost                  float64 `json:"cost"`
	Currency              string  `json:"currency"`
	ScheduledDeliveryDate string  `json:"scheduledDeliveryDate"`
	PhoneNumber           string  `json:"phoneNumber"`
}

type BulkSMSResponse struct {
	Success         bool           `json:"success"`
	Code            int            `json:"code"`
	MessageID       int            `json:"message_id"`
	Currency        string         `json:"currency"`
	Cost            float64        `json:"cost"`
	NumberOfSMS     int            `json:"nbSMS"`
	SMSResponseList []*SMSResponse `json:"SMSResponse_List"`
}

// SendSMS sends SMS, either immediately or at a set time.
func (c *Client) SendSMS(sms *SMS) (*SMSResponse, error) {
	sms.APIKey = c.apiKey
	blob, err := json.Marshal(sms)
	if err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf("%s/send", c.basePath)
	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	blob, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	smsr := new(SMSResponse)
	if err := json.Unmarshal(blob, &smsr); err != nil {
		return nil, err
	}
	return smsr, nil
}

// SendBulkSMS sends SMS in batch of 500 either immediately or at a set time.
func (c *Client) SendBulkSMS(bulksms *BulkSMS) (*BulkSMSResponse, error) {
	bulksms.APIKey = c.apiKey
	blob, err := json.Marshal(bulksms)
	if err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf("%s/bulk-send", c.basePath)
	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	blob, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	bulksmsr := new(BulkSMSResponse)
	if err := json.Unmarshal(blob, bulksmsr); err != nil {
		return nil, err
	}
	return bulksmsr, nil
}

// SendVirtualNumber sends SMS, either immediately or at a set time, with a long number.
func (c *Client) SendVirtualNumber(vn *VNumber) (*SMSResponse, error) {
	vn.APIKey = c.apiKey
	blob, err := json.Marshal(vn)
	if err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf("%s/vn/send", c.basePath)
	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	blob, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	vnr := new(SMSResponse)
	if err := json.Unmarshal(blob, &vnr); err != nil {
		return nil, err
	}
	return vnr, nil
}

type VNumber struct {
	APIKey, To, From, Message string

	// TODO: define optional params
	// IsStopSMS
	// Sandbox
	// Format
}
