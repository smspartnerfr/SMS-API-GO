package smspartner

import "fmt"

// RemoteAPIError is used to handle API error response
// if there are errors (Success == false && Code != 200) the client library
// returns a summary of all errors (e.g., "one error (and 2 other errors)").
// TODO:(hoflish) give client-lib users an option to get verbose errors
type RemoteAPIError struct {
	Success bool               `json:"success,omitempty"`
	Code    int                `json:"code,omitempty"`
	Message string             `json:"message,omitempty"`
	VError  []*ValidationError `json:"error,omitempty"`
}

type ValidationError struct {
	ElementID string `json:"elementId,omitempty"`
	Message   string `json:"message,omitempty"`
}

func (r *RemoteAPIError) Error() string {
	msg, n := "", 0
	for _, e := range r.VError {
		if e != nil {
			if n == 0 {
				msg = e.Message
			}
			n++
		}
	}

	switch n {
	case 0:
		return "(0 errors)"
	case 1:
		return msg
	case 2:
		return msg + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", msg, n-1)
}
