package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		host        *string
		token       *string
		expectError bool
	}{
		{
			name:        "default client",
			host:        nil,
			token:       nil,
			expectError: false,
		},
		{
			name:        "custom host",
			host:        stringPtr("http://custom-host:5678"),
			token:       nil,
			expectError: false,
		},
		{
			name:        "with token",
			host:        nil,
			token:       stringPtr("test-token"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.host, tt.token)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				if tt.host != nil {
					assert.Equal(t, *tt.host, client.HostURL)
				} else {
					assert.Equal(t, HostURL, client.HostURL)
				}
				if tt.token != nil {
					assert.Equal(t, *tt.token, client.Token)
				}
			}
		})
	}
}

func TestDoRequest(t *testing.T) {
	tests := []struct {
		name          string
		server        *httptest.Server
		expectError   bool
		expectMessage string
	}{
		{
			name: "successful request",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"ok"}`))
			})),
			expectError:   false,
			expectMessage: `{"status":"ok"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			host := tt.server.URL
			token := "test"

			client, _ := NewClient(&host, &token)

			req, _ := http.NewRequest("GET", tt.server.URL, nil)
			body, err := client.DoRequest(req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectMessage)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, tt.expectMessage, string(body))
			}
		})
	}
}

func TestGetPaginated(t *testing.T) {
	tests := []struct {
		name          string
		server        *httptest.Server
		expectError   bool
		expectMessage string
	}{
		{
			name: "single page response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[1,2,3]`))
			})),
			expectError:   false,
			expectMessage: `[1,2,3]`,
		},
		{
			name: "paginated response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if r.URL.Query().Get("cursor") == "" {
					w.Write([]byte(`{"data":[1,2],"cursor":"next"}`))
				} else {
					w.Write([]byte(`{"data":[3],"cursor":""}`))
				}
			})),
			expectError:   false,
			expectMessage: `[1,2,3]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			client, _ := NewClient(&host, &token)

			req, _ := http.NewRequest("GET", tt.server.URL, nil)
			body, err := client.GetPaginated(req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, tt.expectMessage, string(body))
			}
		})
	}
}

// Helper function to get a string pointer
func stringPtr(s string) *string {
	return &s
}
