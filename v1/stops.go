package smspartner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DataItem struct {
	ID          int    `json:"id,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

type StopSMSResp struct {
	Success  bool        `json:"success,omitempty"`
	Code     int         `json:"code,omitempty"`
	NbOfData int         `json:"nbData,omitempty"`
	Data     []*DataItem `json:"data,omitempty"`
}

// ListStops returns the list of numbers that sent a STOP.
func (c *Client) ListStops() (*StopSMSResp, error) {
	fullURL := fmt.Sprintf("%s/stop-sms/list?apiKey=%s", c.basePath, c.apiKey)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	str := new(StopSMSResp)
	if err := json.Unmarshal(res, str); err != nil {
		return nil, err
	}
	return str, nil
}

// AddToStops add a phone number to the list of stops.
func (c *Client) AddToStops(phoneNumber string) (map[string]interface{}, error) {
	var payload struct {
		APIKey      string `json:"apiKey,omitempty"`
		PhoneNumber string `json:"phoneNumber,omitempty"`
	}
	payload.APIKey = c.apiKey
	payload.PhoneNumber = phoneNumber
	blob, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s/stop-sms/add", c.basePath)

	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	blob, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(blob, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// DeleteFromStops Deletes a phone number from the list of stops.
func (c *Client) DeleteFromStops(id int) (map[string]interface{}, error) {
	fullURL := fmt.Sprintf("%s/stop-sms/delete?apiKey=%s&id=%d", c.basePath, c.apiKey, id)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(res, &m); err != nil {
		return nil, err
	}
	return m, nil
}
