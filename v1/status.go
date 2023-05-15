package smspartner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SMSStatusResp struct {
	Success     bool    `json:"success,omitempty"`
	Code        int     `json:"code,omitempty"`
	Number      string  `json:"number,omitempty"`
	MessageID   string  `json:"messageId,omitempty"`
	StopSMS     bool    `json:"stopSms,omitempty"`
	Date        string  `json:"date,omitempty"`
	Status      string  `json:"statut,omitempty"`
	Cost        float64 `json:"cost,omitempty"`
	CountryCode string  `json:"countryCode,omitempty"`
	Currency    string  `json:"currency,omitempty"`
	IsSpam      string  `json:"isSpam,omitempty"`
	PhoneNumber string  `json:"phoneNumber,omitempty"`
}

type MultiSMSStatusPayload struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
	MessageID   int    `json:"messageId,omitempty"`
}

type MultiSMSStatusReq struct {
	APIKey        string                   `json:"apiKey,omitempty"`
	SMSStatusList []*MultiSMSStatusPayload `json:"SMSStatut_List,omitempty"`
}

type MultiSMSStatusResp struct {
	Success               bool             `json:"success,omitempty"`
	Code                  int              `json:"code,omitempty"`
	MessageID             string           `json:"message_id,omitempty"`
	SMSStatusResponseList []*SMSStatusResp `json:"StatutResponse_List,omitempty"`
}

// GetSMSStatus returns the status of an SMS
func (c *Client) GetSMSStatus(messageID int, phoneNumber string) (*SMSStatusResp, error) {
	fullURL := fmt.Sprintf("%s/message-status?apiKey=%s&messageId=%d&phoneNumber=%s", c.basePath, c.apiKey, messageID, phoneNumber)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sr := new(SMSStatusResp)
	if err := json.Unmarshal(res, &sr); err != nil {
		return nil, err
	}
	return sr, nil
}

// GetMultiSMSStatus returns the status of multiple SMS
func (c *Client) GetMultiSMSStatus(ss *MultiSMSStatusReq) (*MultiSMSStatusResp, error) {
	ss.APIKey = c.apiKey
	blob, err := json.Marshal(ss)
	if err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s/multi-status", c.basePath)

	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	blob, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	mr := new(MultiSMSStatusResp)
	if err := json.Unmarshal(blob, mr); err != nil {
		return nil, err
	}
	return mr, nil
}

// GetBulkSMSStatus returns the status of multiple SMS by message ID
func (c *Client) GetBulkSMSStatus(messageID int) (*MultiSMSStatusResp, error) {
	fullURL := fmt.Sprintf("%s/bulk-status?apiKey=%s&messageId=%d", c.basePath, c.apiKey, messageID)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	bs := new(MultiSMSStatusResp)
	if err := json.Unmarshal(res, &bs); err != nil {
		return nil, err
	}
	return bs, nil
}
