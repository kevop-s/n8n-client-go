package tags

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kevop-s/n8n-client-go/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestNewTags(t *testing.T) {
	c := &client.Client{}
	tags := NewTags(c)
	assert.NotNil(t, tags)
	assert.Equal(t, c, tags.Client)
}

func TestGetTag(t *testing.T) {
	tests := []struct {
		name         string
		server       *httptest.Server
		tagId        string
		expectedTag  N8nTag
		expectError  bool
		expectStatus int
	}{
		{
			name: "successful get tag",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":        "1",
					"name":      "Test Tag",
					"createdAt": "2023-01-01T00:00:00Z",
					"updatedAt": "2023-01-01T00:00:00Z",
				})
			})),
			tagId: "1",
			expectedTag: N8nTag{
				Id:        "1",
				Name:      "Test Tag",
				CreatedAt: "2023-01-01T00:00:00Z",
				UpdatedAt: "2023-01-01T00:00:00Z",
			},
			expectError:  false,
			expectStatus: http.StatusOK,
		},
		{
			name: "tag not found",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			})),
			tagId:        "999",
			expectError:  true,
			expectStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			tags := NewTags(c)

			tag, err := tags.GetTag(tt.tagId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTag, tag)
			}
		})
	}
}

func TestCreateTag(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		tagName     string
		expectError bool
	}{
		{
			name: "successful create tag",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":        "1",
					"name":      "New Tag",
					"createdAt": "2023-01-01T00:00:00Z",
					"updatedAt": "2023-01-01T00:00:00Z",
				})
			})),
			tagName:     "New Tag",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			tags := NewTags(c)

			tag, err := tags.CreateTag(tt.tagName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.tagName, tag.Name)
			}
		})
	}
}

func TestUpdateTag(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		id          string
		tagName     string
		expectError bool
	}{
		{
			name: "successful update tag",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":        "1",
					"name":      "Updated Tag",
					"createdAt": "2023-01-01T00:00:00Z",
					"updatedAt": "2023-01-02T00:00:00Z",
				})
			})),
			id:          "1",
			tagName:     "Updated Tag",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			tags := NewTags(c)

			tag, err := tags.UpdateTag(tt.id, tt.tagName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.tagName, tag.Name)
			}
		})
	}
}

func TestDeleteTag(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		id          string
		expectError bool
	}{
		{
			name: "successful delete tag",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			id:          "1",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			tags := NewTags(c)

			success, err := tags.DeleteTag(tt.id)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, success)
			}
		})
	}
}
