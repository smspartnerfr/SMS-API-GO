package smspartner

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CancelSMS cancel sending a sent SMS
func (c *Client) CancelSMS(msgID int) (map[string]interface{}, error) {
	fullURL := fmt.Sprintf("%s/message-cancel?apiKey=%s&messageId=%d", c.basePath, c.apiKey, msgID)
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
