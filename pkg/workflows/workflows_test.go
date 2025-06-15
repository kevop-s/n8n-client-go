package workflows

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kevop-s/n8n-client-go/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkflows(t *testing.T) {
	c := &client.Client{}
	w := NewWorkflows(c)
	assert.NotNil(t, w)
	assert.Equal(t, c, w.Client)
}

func TestGetWorkflow(t *testing.T) {
	tests := []struct {
		name           string
		server         *httptest.Server
		expectedError  bool
		expectedResult N8nWorkflow
	}{
		{
			name: "successful get workflow",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":       "1",
					"name":     "Test Workflow",
					"active":   true,
					"nodes":    []interface{}{},
					"settings": map[string]interface{}{},
				})
			})),
			expectedError: false,
			expectedResult: N8nWorkflow{
				Id:     "1",
				Name:   "Test Workflow",
				Active: true,
				Nodes:  []N8nNode{},
				Settings: N8nWorkflowSettings{
					SaveExecutionProgress:    true,
					SaveManualExecutions:     true,
					SaveDataErrorExecution:   "all",
					SaveDataSuccessExecution: "all",
					ExecutionTimeout:         3600,
					ErrorWorkflow:            "",
					Timezone:                 "America/New_York",
					ExecutionOrder:           "v1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"
			c, _ := client.NewClient(&host, &token)
			w := NewWorkflows(c)

			result, err := w.GetWorkflow("1")

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Id, result.Id)
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Active, result.Active)
			}
		})
	}
}

func TestCreateWorkflow(t *testing.T) {
	tests := []struct {
		name         string
		server       *httptest.Server
		workflowData N8nWorkflow
		expectError  bool
	}{
		{
			name: "successful workflow creation",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":          "new-workflow",
					"name":        "New Workflow",
					"connections": map[string]interface{}{},
				})
			})),
			workflowData: N8nWorkflow{
				Name: "New Workflow",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"
			c, _ := client.NewClient(&host, &token)
			w := NewWorkflows(c)

			result, err := w.CreateWorkflow(tt.workflowData)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "new-workflow", result.Id)
				assert.Equal(t, "New Workflow", result.Name)
			}
		})
	}
}

func TestUpdateWorkflow(t *testing.T) {
	tests := []struct {
		name         string
		server       *httptest.Server
		workflowData N8nWorkflow
		expectError  bool
	}{
		{
			name: "successful workflow update",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   "1",
					"name": "Updated Workflow",
				})
			})),
			workflowData: N8nWorkflow{
				Name: "Updated Workflow",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"
			c, _ := client.NewClient(&host, &token)
			w := NewWorkflows(c)

			result, err := w.UpdateWorkflow("1", tt.workflowData)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "1", result.Id)
				assert.Equal(t, "Updated Workflow", result.Name)
			}
		})
	}
}

func TestDeleteWorkflow(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		expectError bool
	}{
		{
			name: "successful workflow deletion",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"
			c, _ := client.NewClient(&host, &token)
			w := NewWorkflows(c)

			success, err := w.DeleteWorkflow("1")
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, success)
			}
		})
	}
}
