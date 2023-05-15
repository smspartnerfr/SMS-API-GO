package smspartner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const envSMSPartnerAPIKey = "SMSPARTNER_API_KEY"
const apiBasePath = "http://api.smspartner.fr/v1"
const clientDefaultTimeout time.Duration = 10 * time.Second

var errUnsetAPIKey = fmt.Errorf("could not find %q in your environment", envSMSPartnerAPIKey)

type Client struct {
	hc       *http.Client
	basePath string
	apiKey   string
}

// NewClient returns a HTTP client.
func NewClient(c *http.Client, opts ...Option) (*Client, error) {
	wrapClient := new(http.Client)
	*wrapClient = *c

	t := c.Timeout
	if t == 0 {
		t = clientDefaultTimeout
	}
	tr := c.Transport
	if tr == nil {
		tr = http.DefaultTransport
	}

	wrapClient.Timeout = t
	wrapClient.Transport = tr

	apiKey, err := getAPIKeyFromEnv()
	if err != nil {
		return nil, err
	}

	client := &Client{
		hc:       wrapClient,
		apiKey:   apiKey,
		basePath: apiBasePath,
	}

	if err := client.parseOptions(opts...); err != nil {
		return nil, err
	}

	return client, nil
}

func getAPIKeyFromEnv() (string, error) {
	apikey := strings.TrimSpace(os.Getenv(envSMSPartnerAPIKey))
	if apikey == "" {
		return "", errUnsetAPIKey
	}
	return apikey, nil
}

type Option func(*Client) error

func BasePath(basePath string) Option {
	return func(c *Client) error {
		c.basePath = basePath
		return nil
	}
}

func APIKey(apiKey string) Option {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}

func (c *Client) parseOptions(opts ...Option) error {
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response: %v", err)
	}

	// handle non-200 status code
	if resp.StatusCode != http.StatusOK {
		remAPIErr := &RemoteAPIError{}
		if err := json.Unmarshal(body, remAPIErr); err != nil {
			return nil, fmt.Errorf("error unmarshalling response: %v", err)
		}

		if !remAPIErr.Success && remAPIErr.Code != 200 {
			if remAPIErr.Message != "" {
				return nil, errors.New(remAPIErr.Message)
			}
			return nil, errors.New(remAPIErr.Error())
		}

		if remAPIErr == nil {
			return nil, fmt.Errorf("unexpected response: %s", string(body))
		}
	}
	return body, nil
}
