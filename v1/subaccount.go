package smspartner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type SubAccountType string

const (
	Simple   SubAccountType = "simple"
	Advanced SubAccountType = "advanced"
)

type SubAccountCreationRequest struct {
	APIKey     string                    `json:"apiKey,omitempty"`
	Type       SubAccountType            `json:"type,omitempty"`
	Parameters *SubAccountCreationParams `json:"parameters,omitempty"`
}

type SubAccountCreationParams struct {
	Email             string `json:"email,omitempty"`
	IsBuyer           int    `json:"isBuyer,omitempty"`
	CreditToAttribute int    `json:"creditToAttribute,omitempty"`
	FirstName         string `json:"firstname,omitempty"`
	LastName          string `json:"lastname,omitempty"`
}

type SubAccountCreationResponse struct {
	Success    bool   `json:"success"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	SubAccount struct {
		Email string `json:"email"`
		Token string `json:"token"`
	} `json:"subaccount"`
	SendConfirmMailTo string `json:"sendConfirmMailTo"`
	ParentEmail       string `json:"parent_email"`
}

var ErrSubAccountEmail = errors.New("Email is required for Advanced sub-account")

// CreateSubAccount creates a sub account.
func (c *Client) CreateSubAccount(subAccReq *SubAccountCreationRequest) (*SubAccountCreationResponse, error) {
	if subAccReq.Type == Advanced && subAccReq.Parameters != nil {
		if subAccReq.Parameters.Email == "" {
			return nil, ErrSubAccountEmail
		}
	}
	subAccReq.APIKey = c.apiKey

	blob, err := json.Marshal(subAccReq)
	if err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s/subaccount/create", c.basePath)

	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	blob, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	subAccResp := new(SubAccountCreationResponse)
	if err := json.Unmarshal(blob, &subAccResp); err != nil {
		return nil, err
	}
	return subAccResp, nil
}

type SubAccount struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Type      string `json:"type"`
	Token     string `json:"token"`
	APIKey    string `json:"apiKey"`
	CreatedAt string `json:"createdAt"`
	Credits   struct {
		Balance  string `json:"balance"`
		Currency string `json:"currency"`
	} `json:"credits"`
}

type SubAccountsResponse struct {
	Success   bool          `json:"success"`
	Code      int           `json:"code"`
	Message   string        `json:"message"`
	Total     int           `json:"total"`
	NbPerPage int           `json:"nb_per_page"`
	Page      int           `json:"page"`
	Data      []*SubAccount `json:"data"`
}

// ListSubAccounts lists all sub accounts.
func (c *Client) ListSubAccounts() (*SubAccountsResponse, error) {
	fullURL := fmt.Sprintf("%s/subaccount/list?apiKey=%s", c.basePath, c.apiKey)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	subAccsResp := new(SubAccountsResponse)
	if err := json.Unmarshal(res, subAccsResp); err != nil {
		return nil, err
	}
	return subAccsResp, nil
}

type SubAccountCreditAdditionResponse struct {
	Success          bool    `json:"success"`
	Code             int     `json:"code"`
	Message          string  `json:"message"`
	Credit           float64 `json:"total"`
	SubaccountCredit float64 `json:"subaccountCredit"`
	Currency         string  `json:"currency"`
}

// AddCreditToSubAccount - Credits will be debited from the main account.
func (c *Client) AddCreditToSubAccount(credit, tokenSubaccount string) (*SubAccountCreditAdditionResponse, error) {
	var payload struct {
		APIKey          string `json:"apiKey,omitempty"`
		Credit          string `json:"credit,omitempty"`
		TokenSubAccount string `json:"tokenSubaccount,omitempty"`
	}
	payload.APIKey = c.apiKey
	payload.Credit = credit
	payload.TokenSubAccount = tokenSubaccount

	blob, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s/subaccount/credit/add", c.basePath)

	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	blob, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sr := new(SubAccountCreditAdditionResponse)
	if err := json.Unmarshal(blob, &sr); err != nil {
		return nil, err
	}
	return sr, nil
}

// DeleteCreditFromSubAccount -
// REVIEW: https://my.smspartner.fr/documentation-fr/api/v1/subaccount/credit/remove
func (c *Client) DeleteCreditFromSubAccount() {
	panic("Not implemented yet!")
}
