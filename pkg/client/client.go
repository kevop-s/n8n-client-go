package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default Hashicups URL
const HostURL string = "http://localhost:19090"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient -
func NewClient(host, token *string) (*Client, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second, Transport: tr},
		HostURL:    HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if token != nil {
		c.Token = *token
	}

	return &c, nil
}

// N8nPaginatedResponse represents a paginated API response with a cursor and data
type N8nPaginatedResponse struct {
	Data   json.RawMessage `json:"data"`
	Cursor string          `json:"cursor,omitempty"`
}

type N8nErrorResponse struct {
	Message string `json:"message"`
}

// DoRequest performs an HTTP request and returns the response body
func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("no token provided")
	}

	req.Header.Set("X-N8n-API-KEY", c.Token)
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var errorResp N8nErrorResponse
		json.Unmarshal(body, &errorResp)
		return nil, fmt.Errorf("status: %d, message: %s", res.StatusCode, errorResp.Message)
	}

	return body, err
}

// GetPaginated performs a GET request with cursor-based pagination support
// It automatically handles fetching all pages
func (c *Client) GetPaginated(req *http.Request) ([]byte, error) {
	var allData []json.RawMessage
	cursor := ""

	for {
		// Create URL with query parameters
		if cursor != "" {
			q := req.URL.Query()
			q.Set("cursor", cursor)
			req.URL.RawQuery = q.Encode()
		}

		body, err := c.DoRequest(req)
		if err != nil {
			return nil, err
		}

		// Try to unmarshal as paginated response
		var page N8nPaginatedResponse
		if err := json.Unmarshal(body, &page); err != nil || len(page.Data) == 0 {
			// If not a paginated response, return as is
			if len(allData) == 0 {
				return body, nil
			}
			// If we already have data, append the current response
			allData = append(allData, body)
			break
		}

		// Append current page data
		var currentData []json.RawMessage
		if err := json.Unmarshal(page.Data, &currentData); err != nil {
			// If data is not an array, add it as a single item
			if len(page.Data) > 0 {
				allData = append(allData, page.Data)
			}
		} else {
			allData = append(allData, currentData...)
		}

		// If no more pages, break
		if page.Cursor == "" {
			break
		}
		cursor = page.Cursor
	}

	// Combine all data into a single JSON array
	result, err := json.Marshal(allData)
	if err != nil {
		return nil, fmt.Errorf("error combining results: %v", err)
	}

	return result, nil
}
