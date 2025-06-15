package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kevop-s/n8n-client-go/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestNewUsers(t *testing.T) {
	c := &client.Client{}
	users := NewUsers(c)
	assert.NotNil(t, users)
	assert.Equal(t, c, users.Client)
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name         string
		server       *httptest.Server
		userId       string
		expectedUser N8nUser
		expectError  bool
	}{
		{
			name: "successful get user",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":        "1",
					"email":     "test@example.com",
					"firstName": "Test",
					"lastName":  "User",
					"isPending": false,
					"role":      "admin",
				})
			})),
			userId: "1",
			expectedUser: N8nUser{
				Id:        "1",
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
				IsPending: false,
				Role:      "admin",
			},
			expectError: false,
		},
		{
			name: "user not found",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			})),
			userId:      "999",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			users := NewUsers(c)

			user, err := users.GetUser(tt.userId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		email       string
		role        string
		expectError bool
	}{
		{
			name: "successful create user",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":        "2",
					"email":     "new@example.com",
					"firstName": "",
					"lastName":  "",
					"isPending": true,
					"role":      "member",
				})
			})),
			email:       "new@example.com",
			role:        "member",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			users := NewUsers(c)

			user, err := users.CreateUser(tt.email, tt.role)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.role, user.Role)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		email       string
		role        string
		expectError bool
	}{
		{
			name: "successful update user",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":        "1",
					"email":     "updated@example.com",
					"firstName": "Updated",
					"lastName":  "User",
					"isPending": false,
					"role":      "admin",
				})
			})),
			email:       "updated@example.com",
			role:        "admin",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			users := NewUsers(c)

			user, err := users.UpdateUser(tt.email, tt.role)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.role, user.Role)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		email       string
		expectError bool
	}{
		{
			name: "successful delete user",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			email:       "user@example.com",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			users := NewUsers(c)

			success, err := users.DeleteUser(tt.email)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, success)
			}
		})
	}
}

// Helper function to get a string pointer
func stringPtr(s string) *string {
	return &s
}
