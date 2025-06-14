package workflows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevop-s/n8n-client-go/pkg/client"
)

type Workflows struct {
	Client *client.Client
}

type N8nWorkflow struct {
	Id          string                   `json:"id,omitempty"`
	Name        string                   `json:"name"`
	Active      bool                     `json:"active,omitempty"`
	Nodes       []N8nNode                `json:"nodes"`
	Connections map[string][]interface{} `json:"connections"`
	Settings    N8nWorkflowSettings      `json:"settings"`
	StaticData  string                   `json:"staticData,omitempty"`
}

type N8nWorkflowSettings struct {
	SaveExecutionProgress    bool   `json:"saveExecutionProgress,omitempty"`
	SaveManualExecutions     bool   `json:"saveManualExecutions,omitempty"`
	SaveDataErrorExecution   string `json:"saveDataErrorExecution,omitempty"`
	SaveDataSuccessExecution string `json:"saveDataSuccessExecution,omitempty"`
	ExecutionTimeout         int    `json:"executionTimeout,omitempty"`
	ErrorWorkflow            string `json:"errorWorkflow,omitempty"`
	Timezone                 string `json:"timezone,omitempty"`
	ExecutionOrder           string `json:"executionOrder,omitempty"`
}

func NewWorkflows(client *client.Client) *Workflows {
	return &Workflows{Client: client}
}

// GetWorkflow retrieves a workflow by its ID
func (w *Workflows) GetWorkflow(id string) (N8nWorkflow, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/workflows/%s", w.Client.HostURL, id), nil)
	if err != nil {
		return N8nWorkflow{}, err
	}
	resp, err := w.Client.GetPaginated(req)

	if err != nil {
		return N8nWorkflow{}, err
	}

	var workflow N8nWorkflow
	err = json.Unmarshal(resp, &workflow)

	if err != nil {
		return N8nWorkflow{}, err
	}

	return workflow, nil
}

// CreateWorkflow creates a new workflow
func (w *Workflows) CreateWorkflow(workflowData N8nWorkflow) (N8nWorkflow, error) {
	if len(workflowData.Nodes) > 0 {
		return N8nWorkflow{}, fmt.Errorf("nodes can not be defined when workflow is created, use AddNode instead")
	}

	if workflowData.Connections != nil {
		return N8nWorkflow{}, fmt.Errorf("connections can not be defined when workflow is created, use AddConnection instead")
	}

	workflowData.Nodes = []N8nNode{}
	workflowData.Connections = map[string][]interface{}{}
	w.setDefaultWorkflowSettings(&workflowData)

	jsonWorkflow, err := json.Marshal(workflowData)
	if err != nil {
		return N8nWorkflow{}, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/workflows", w.Client.HostURL), bytes.NewReader(jsonWorkflow))
	if err != nil {
		return N8nWorkflow{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.Client.DoRequest(req)

	if err != nil {
		return N8nWorkflow{}, err
	}
	var workflow N8nWorkflow
	err = json.Unmarshal(resp, &workflow)

	if err != nil {
		return N8nWorkflow{}, err
	}

	return workflow, nil
}

// UpdateWorkflow updates an existing workflow
func (w *Workflows) UpdateWorkflow(id string, workflowData N8nWorkflow) (N8nWorkflow, error) {
	currentWorkflow, err := w.GetWorkflow(id)
	if err != nil {
		return N8nWorkflow{}, err
	}

	combinedWorkflowData := w.combineWorkflows(currentWorkflow, workflowData)
	// remove readonly fields
	combinedWorkflowData.Id = ""
	combinedWorkflowData.Nodes = currentWorkflow.Nodes
	combinedWorkflowData.Connections = currentWorkflow.Connections

	if len(workflowData.Nodes) > 0 {
		combinedWorkflowData.Nodes = workflowData.Nodes
	}

	if len(workflowData.Connections) > 0 {
		combinedWorkflowData.Connections = workflowData.Connections
	}

	jsonWorkflow, err := json.Marshal(combinedWorkflowData)
	if err != nil {
		return N8nWorkflow{}, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/workflows/%s", w.Client.HostURL, id), bytes.NewReader(jsonWorkflow))
	if err != nil {
		return N8nWorkflow{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := w.Client.DoRequest(req)

	if err != nil {
		return N8nWorkflow{}, err
	}
	var workflow N8nWorkflow
	err = json.Unmarshal(resp, &workflow)

	if err != nil {
		return N8nWorkflow{}, err
	}

	return workflow, nil
}

// DeleteWorkflow deletes a workflow by its ID
func (w *Workflows) DeleteWorkflow(id string) (bool, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/workflows/%s", w.Client.HostURL, id), nil)
	if err != nil {
		return false, err
	}
	_, err = w.Client.DoRequest(req)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (w *Workflows) combineWorkflows(originalWorkflow N8nWorkflow, updateWorkflow N8nWorkflow) N8nWorkflow {

	jsonWorkflowOriginal, err := json.Marshal(originalWorkflow)
	if err != nil {
		return N8nWorkflow{}
	}
	jsonWorkflowUpdate, err := json.Marshal(updateWorkflow)
	if err != nil {
		return N8nWorkflow{}
	}

	var finalWorkflow N8nWorkflow
	err = json.Unmarshal(jsonWorkflowOriginal, &finalWorkflow)
	if err != nil {
		return N8nWorkflow{}
	}

	err = json.Unmarshal(jsonWorkflowUpdate, &finalWorkflow)
	if err != nil {
		return N8nWorkflow{}
	}

	return finalWorkflow
}

func (w *Workflows) setDefaultWorkflowSettings(workflow *N8nWorkflow) {
	if workflow.Settings.SaveDataErrorExecution == "" {
		workflow.Settings.SaveDataErrorExecution = "all"
	}
	if workflow.Settings.SaveDataSuccessExecution == "" {
		workflow.Settings.SaveDataSuccessExecution = "all"
	}
	if workflow.Settings.ExecutionTimeout == 0 {
		workflow.Settings.ExecutionTimeout = 3600
	}
	if workflow.Settings.ErrorWorkflow == "" {
		workflow.Settings.ErrorWorkflow = "all"
	}
	if workflow.Settings.Timezone == "" {
		workflow.Settings.Timezone = "UTC"
	}
	if workflow.Settings.ExecutionOrder == "" {
		workflow.Settings.ExecutionOrder = "v1"
	}
}
