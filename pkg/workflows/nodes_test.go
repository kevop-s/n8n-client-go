package workflows

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kevop-s/n8n-client-go/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestGetNodes(t *testing.T) {
	tests := []struct {
		name           string
		server         *httptest.Server
		expectedError  bool
		expectedLength int
	}{
		{
			name: "successful get nodes",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   "1",
					"name": "Test Workflow",
					"nodes": []map[string]interface{}{
						{
							"parameters": map[string]interface{}{
								"options": map[string]interface{}{},
							},
							"id":               "b24b05a7-d802-4413-bfb1-23e1e76f6203",
							"name":             "When chat message received",
							"webhookId":        "a889d2ae-2159-402f-b326-5f61e90f602e",
							"notesInFlow":      false,
							"type":             "@n8n/n8n-nodes-langchain.chatTrigger",
							"typeVersion":      1.1,
							"executeOnce":      false,
							"alwaysOutputData": false,
							"retryOnFail":      false,
							"maxTries":         0,
							"waitBetweenTries": 0,
							"position":         []int{360, 20},
							"onError":          "",
						},
					},
					"connections": map[string]interface{}{},
					"settings": map[string]interface{}{
						"executionOrder": "v1",
					},
				})
			})),
			expectedError:  false,
			expectedLength: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			w := NewWorkflows(c)

			nodes, err := w.GetNodes("1")
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLength, len(nodes))
			}
		})
	}
}

func TestGetNodeById(t *testing.T) {
	tests := []struct {
		name         string
		server       *httptest.Server
		nodeId       string
		expectError  bool
		expectedName string
	}{
		{
			name: "successful get node by id",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   "1",
					"name": "Test Workflow",
					"nodes": []map[string]interface{}{
						{
							"parameters": map[string]interface{}{
								"options": map[string]interface{}{},
							},
							"id":               "b24b05a7-d802-4413-bfb1-23e1e76f6203",
							"name":             "When chat message received",
							"webhookId":        "a889d2ae-2159-402f-b326-5f61e90f602e",
							"notesInFlow":      false,
							"type":             "@n8n/n8n-nodes-langchain.chatTrigger",
							"typeVersion":      1.1,
							"executeOnce":      false,
							"alwaysOutputData": false,
							"retryOnFail":      false,
							"maxTries":         0,
							"waitBetweenTries": 0,
							"position":         []int{360, 20},
							"onError":          "",
						},
					},
					"connections": map[string]interface{}{},
					"settings": map[string]interface{}{
						"executionOrder": "v1",
					},
				})
			})),
			nodeId:       "b24b05a7-d802-4413-bfb1-23e1e76f6203",
			expectError:  false,
			expectedName: "When chat message received",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			host := tt.server.URL
			token := "test"

			c, _ := client.NewClient(&host, &token)
			w := NewWorkflows(c)

			node, err := w.GetNodeById("1", tt.nodeId)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedName, node.Name)
			}
		})
	}
}

func TestAddNode(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		node        N8nNode
		expectError bool
	}{
		{
			name: "successful add node",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"id":          "1",
						"name":        "Test Workflow",
						"nodes":       []map[string]interface{}{},
						"connections": map[string]interface{}{},
						"settings": map[string]interface{}{
							"executionOrder": "v1",
						},
					})
				} else {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"id":   "1",
						"name": "Test Workflow",
						"nodes": []map[string]interface{}{
							{
								"parameters": map[string]interface{}{
									"options": map[string]interface{}{},
								},
								"id":               "b24b05a7-d802-4413-bfb1-23e1e76f6203",
								"name":             "When chat message received",
								"webhookId":        "a889d2ae-2159-402f-b326-5f61e90f602e",
								"notesInFlow":      false,
								"type":             "@n8n/n8n-nodes-langchain.chatTrigger",
								"typeVersion":      1.1,
								"executeOnce":      false,
								"alwaysOutputData": false,
								"retryOnFail":      false,
								"maxTries":         0,
								"waitBetweenTries": 0,
								"position":         []int{100, 100},
								"onError":          "",
							},
						},
						"connections": map[string]interface{}{},
						"settings": map[string]interface{}{
							"executionOrder": "v1",
						},
					})
				}
			})),
			node: N8nNode{
				Name:     "When chat message received",
				Type:     "@n8n/n8n-nodes-langchain.chatTrigger",
				Position: []int{100, 100},
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

			_, err := w.AddNode("1", tt.node)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateNode(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		nodeId      string
		updateNode  N8nNode
		expectError bool
	}{
		{
			name: "successful update node",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"id":   "1",
						"name": "Test Workflow",
						"nodes": []map[string]interface{}{
							{"id": "node1", "name": "Old Name"},
						},
					})
				} else {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"id":   "1",
						"name": "Test Workflow",
					})
				}
			})),
			nodeId: "node1",
			updateNode: N8nNode{
				Name:     "Updated Node",
				Type:     "n8n-nodes-base.httpRequest",
				Position: []int{100, 100},
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

			_, err := w.UpdateNode("1", tt.nodeId, tt.updateNode)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveNode(t *testing.T) {
	tests := []struct {
		name        string
		server      *httptest.Server
		nodeId      string
		expectError bool
	}{
		{
			name: "successful remove node",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"id":   "1",
						"name": "Test Workflow",
						"nodes": []map[string]interface{}{
							{"id": "node1", "name": "Node 1"},
						},
						"connections": map[string]interface{}{},
						"settings": map[string]interface{}{
							"executionOrder": "v1",
						},
					})
				} else {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"id":          "1",
						"name":        "Test Workflow",
						"nodes":       []map[string]interface{}{},
						"connections": map[string]interface{}{},
						"settings": map[string]interface{}{
							"executionOrder": "v1",
						},
					})
				}
			})),
			nodeId:      "node1",
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

			success, err := w.RemoveNode("1", tt.nodeId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, success)
			}
		})
	}
}
