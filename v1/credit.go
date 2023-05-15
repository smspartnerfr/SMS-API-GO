package smspartner

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Username  string `json:"username"`
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type Credits struct {
	Balance          string `json:"balance"`
	CreditHlr        int    `json:"creditHlr"`
	CreditSMS        int    `json:"creditSms"`
	CreditSmsLowCost int    `json:"creditSmsLowCost"`
	Currency         string `json:"currency"`
	ToSend           int    `json:"toSend"`
	Solde            string `json:"solde"`
}

type CreditsResponse struct {
	Success bool     `json:"success"`
	Code    int      `json:"code"`
	User    *User    `json:"user"`
	Credits *Credits `json:"credits"`
}

// CheckCredits returns your SMS credit (number of SMS available, based on your
// own purchases and usage), as well as the number of SMS that are about to be sent.
func (c *Client) CheckCredits() (*CreditsResponse, error) {
	fullURL := fmt.Sprintf("%s/me?apiKey=%s", c.basePath, c.apiKey)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	credits := new(CreditsResponse)
	if err := json.Unmarshal(res, credits); err != nil {
		return nil, err
	}
	return credits, nil
}
