package workflows

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kevop-s/n8n-client-go/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestGetConnectionBySourceNodeName(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		workflowId  string
		nodeName    string
		expected    N8nConnection
		expectError bool
	}{
		{
			name: "successful get connection by source node name",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   "1",
					"name": "Test Workflow",
					"connections": map[string]interface{}{
						"When chat message received": map[string]interface{}{
							"main": []interface{}{
								[]interface{}{
									map[string]interface{}{
										"node":  "Agent",
										"type":  "main",
										"index": 0,
									},
								},
							},
						},
					},
				})
			})),
			workflowId: "1",
			nodeName:   "When chat message received",
			expected: N8nConnection{
				SourceNodeName: "When chat message received",
				ConnectionType: "main",
				Outputs: []N8nConnectionOutput{
					{
						OutputIndex:               0,
						DestinationNodeName:       "Agent",
						DestinationNodeInputType:  "main",
						DestinationNodeInputIndex: 0,
					},
				},
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
			workflows := NewWorkflows(c)

			connection, err := workflows.GetConnectionBySourceNodeName(tt.workflowId, tt.nodeName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, connection)
			}
		})
	}
}

func TestGetConnections(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		workflowId  string
		expected    []N8nConnection
		expectError bool
	}{
		{
			name: "successful get connections",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   "1",
					"name": "Test Workflow",
					"connections": map[string]interface{}{
						"When chat message received": map[string]interface{}{
							"main": []interface{}{
								[]interface{}{
									map[string]interface{}{
										"node":  "Agent",
										"type":  "main",
										"index": 0,
									},
								},
							},
						},
					},
				})
			})),
			workflowId: "1",
			expected: []N8nConnection{
				{
					SourceNodeName: "When chat message received",
					ConnectionType: "main",
					Outputs: []N8nConnectionOutput{
						{
							OutputIndex:               0,
							DestinationNodeName:       "Agent",
							DestinationNodeInputIndex: 0,
							DestinationNodeInputType:  "main",
						},
					},
				},
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
			workflows := NewWorkflows(c)

			connections, err := workflows.GetConnections(tt.workflowId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, connections)
			}
		})
	}
}

func TestAddConnection(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		workflowId  string
		connection  N8nConnection
		expect      N8nConnection
		expectError bool
	}{
		{
			name: "successful add connection",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"id":          "1",
						"name":        "Test Workflow",
						"connections": map[string]interface{}{},
					})
				} else {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"id":   "1",
						"name": "Test Workflow",
						"connections": map[string]interface{}{
							"When chat message received": map[string]interface{}{
								"main": []interface{}{
									[]interface{}{
										map[string]interface{}{
											"node":  "Agent",
											"type":  "main",
											"index": 0,
										},
									},
								},
							},
						},
					})
				}
			})),
			workflowId: "1",
			connection: N8nConnection{
				SourceNodeName: "When chat message received",
				ConnectionType: "main",
				Outputs: []N8nConnectionOutput{
					{
						OutputIndex:               0,
						DestinationNodeName:       "Agent",
						DestinationNodeInputIndex: 0,
						DestinationNodeInputType:  "main",
					},
				},
			},
			expect: N8nConnection{
				SourceNodeName: "When chat message received",
				ConnectionType: "main",
				Outputs: []N8nConnectionOutput{
					{
						OutputIndex:               0,
						DestinationNodeName:       "Agent",
						DestinationNodeInputIndex: 0,
						DestinationNodeInputType:  "main",
					},
				},
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
			workflows := NewWorkflows(c)

			connection, err := workflows.AddConnection(tt.workflowId, tt.connection)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, connection)
			}
		})
	}
}

func TestRemoveConnection(t *testing.T) {
	tests := []struct {
		name           string
		server         *httptest.Server
		workflowId     string
		sourceNodeName string
		expectError    bool
	}{
		{
			name: "successful remove connection",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":          "1",
					"name":        "Test Workflow",
					"connections": map[string]interface{}{},
				})
				w.WriteHeader(http.StatusOK)
			})),
			workflowId:     "1",
			sourceNodeName: "When chat message received",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			workflows := NewWorkflows(c)

			success, err := workflows.RemoveConnection(tt.workflowId, tt.sourceNodeName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, success)
			}
		})
	}
}

func TestUpdateConnection(t *testing.T) {
	tests := []struct {
		name           string
		server         *httptest.Server
		workflowId     string
		sourceNodeName string
		connection     N8nConnection
		expect         N8nConnection
		expectError    bool
	}{
		{
			name: "successful update connection",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   "1",
					"name": "Test Workflow",
					"connections": map[string]interface{}{
						"When chat message received": map[string]interface{}{
							"main": []interface{}{
								[]interface{}{
									map[string]interface{}{
										"node":  "Agent",
										"type":  "main",
										"index": 0,
									},
								},
							},
						},
					},
				})
			})),
			workflowId:     "1",
			sourceNodeName: "When chat message received",
			connection: N8nConnection{
				SourceNodeName: "When chat message received",
				ConnectionType: "updated",
				Outputs: []N8nConnectionOutput{
					{
						OutputIndex:               0,
						DestinationNodeName:       "Agent",
						DestinationNodeInputIndex: 0,
						DestinationNodeInputType:  "main",
					},
				},
			},
			expect: N8nConnection{
				SourceNodeName: "When chat message received",
				ConnectionType: "updated",
				Outputs: []N8nConnectionOutput{
					{
						OutputIndex:               0,
						DestinationNodeName:       "Agent",
						DestinationNodeInputIndex: 0,
						DestinationNodeInputType:  "main",
					},
				},
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
			workflows := NewWorkflows(c)

			connection, err := workflows.UpdateConnection(tt.workflowId, tt.sourceNodeName, tt.connection)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, connection)
			}
		})
	}
}
