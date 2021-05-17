package paystack

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	// defaultHTTPTimeout is the default timeout on the http client
	defaultHTTPTimeout = 60 * time.Second

	// base URL for all Paystack API requests
	baseURL = "https://api.paystack.co"

	// User agent used when communicating with the Paystack API.
	// userAgent = "paystack-go/" + version
	userAgent = "Mozilla/5.0 (Unknown; Linux) AppleWebKit/538.1 (KHTML, like Gecko) Chrome/v1.0.0 Safari/538.1"
)

type service struct {
	client *Client
}

type Client struct {
	common service      // Reuse a single struct instead of allocating one for each service on the heap.
	client *http.Client // HTTP client used to communicate with the API.
	// the API Key used to authenticate all Paystack API requests
	key     string
	baseURL *url.URL

	Transaction *TransactionService
}

// Metadata is an key-value pairs added to Paystack API requests
type Metadata map[string]interface{}

// Response represents arbitrary response data
type Response map[string]interface{}

// ListMeta is pagination metadata for paginated responses from the Paystack API
type ListMeta struct {
	Total     int `json:"total"`
	Skipped   int `json:"skipped"`
	PerPage   int `json:"perPage"`
	Page      int `json:"page"`
	PageCount int `json:"pageCount"`
}

// RequestValues aliased to url.Values as a workaround
type RequestValues url.Values

// NewClient creates a new Paystack API client with the given API key
// and HTTP client, allowing overriding of the HTTP client to use.
// This is useful if you're running in a Google AppEngine environment
// where the http.DefaultClient is not available.
func NewClient(key string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}

	u, _ := url.Parse(baseURL)

	c := &Client{
		client:  httpClient,
		key:     key,
		baseURL: u,
	}

	c.common.client = c
	c.Transaction = (*TransactionService)(&c.common)
	return c
}

// Call actually does the HTTP request to Paystack API
func (c *Client) Call(method, path string, body, v interface{}) error {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return err
		}
	}
	u, _ := c.baseURL.Parse(path)
	req, err := http.NewRequest(method, u.String(), buf)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return c.decodeResponse(resp, v)
}

func mapstruct(data interface{}, v interface{}) error {
	config := &mapstructure.DecoderConfig{
		Result:           v,
		TagName:          "json",
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	err = decoder.Decode(data)
	return err
}

// decodeResponse decodes the JSON response from the payment API.
// The actual response will be written to the `v` parameter
func (c *Client) decodeResponse(httpResp *http.Response, v interface{}) error {
	var resp Response
	respBody, _ := ioutil.ReadAll(httpResp.Body)
	json.Unmarshal(respBody, &resp)

	if status, _ := resp["status"].(bool); !status || httpResp.StatusCode >= 400 {

		return newAPIError(httpResp)
	}

	if data, ok := resp["data"]; ok {
		switch t := resp["data"].(type) {
		case map[string]interface{}:
			return mapstruct(data, v)
		default:
			_ = t
			return mapstruct(resp, v)
		}
	}
	// if response data does not contain data key, map entire response to v
	return mapstruct(resp, v)
}

type InitializePaymentType struct {
	Amount    float64
	Email     string
	reference string
	metadata  *map[string]interface{}
}

func initializePayment(payment InitializePaymentType) {

}
