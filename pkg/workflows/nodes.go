package workflows

import (
	"encoding/json"
	"fmt"
)

type N8nNode struct {
	Id               string                 `json:"id,omitempty"`
	Name             string                 `json:"name"`
	WebhookId        string                 `json:"webhookId"`
	Disabled         bool                   `json:"disabled"`
	NotesInFlow      bool                   `json:"notesInFlow"`
	Notes            string                 `json:"notes"`
	Type             string                 `json:"type"`
	TypeVersion      float64                `json:"typeVersion"`
	ExecuteOnce      bool                   `json:"executeOnce"`
	AlwaysOutputData bool                   `json:"alwaysOutputData"`
	RetryOnFail      bool                   `json:"retryOnFail"`
	MaxTries         int                    `json:"maxTries"`
	WaitBetweenTries int                    `json:"waitBetweenTries"`
	ContinueOnFail   bool                   `json:"continueOnFail"`
	OnError          string                 `json:"onError"`
	Parameters       map[string]interface{} `json:"parameters,omitempty"`
	Position         []int                  `json:"position"`
	Credentials      map[string]interface{} `json:"credentials,omitempty"`
}

// GetNodes retrieves all nodes from a workflow
func (w *Workflows) GetNodes(workflowId string) ([]N8nNode, error) {
	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return nil, err
	}

	return workflow.Nodes, nil
}

// GetNodeById retrieves a node by its ID
func (w *Workflows) GetNodeById(workflowId string, nodeId string) (N8nNode, error) {
	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return N8nNode{}, err
	}

	for _, node := range workflow.Nodes {
		if node.Id == nodeId {
			return node, nil
		}
	}

	return N8nNode{}, nil
}

// GetNodeByName retrieves a node by its name
func (w *Workflows) GetNodeByName(workflowId string, nodeName string) (N8nNode, error) {
	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return N8nNode{}, err
	}

	for _, node := range workflow.Nodes {
		if node.Name == nodeName {
			return node, nil
		}
	}

	return N8nNode{}, nil
}

// AddNode adds a new node to a workflow
func (w *Workflows) AddNode(workflowId string, newNode N8nNode) (N8nNode, error) {
	if err := w.validateNodeInput(newNode); err != nil {
		return N8nNode{}, err
	}

	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return N8nNode{}, err
	}

	for _, node := range workflow.Nodes {
		if node.Name == newNode.Name {
			return N8nNode{}, fmt.Errorf("node already exists, use a different name")
		}
	}

	var finalNodes []N8nNode

	finalNodes = append(finalNodes, newNode)
	// remove readonly fields
	newNode.Id = ""

	for _, node := range workflow.Nodes {
		// remove readonly fields
		newNode.Id = ""
		finalNodes = append(finalNodes, node)
	}

	workflow.Nodes = finalNodes

	_, err = w.UpdateWorkflow(workflowId, workflow)

	if err != nil {
		return N8nNode{}, err
	}

	addedNode, err := w.GetNodeByName(workflowId, newNode.Name)

	if err != nil {
		return N8nNode{}, err
	}

	return addedNode, nil

}

// RemoveNode removes a node from a workflow
func (w *Workflows) RemoveNode(workflowId string, nodeId string) (bool, error) {
	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return false, err
	}

	var finalNodes []N8nNode

	for _, node := range workflow.Nodes {
		if node.Id != nodeId {
			finalNodes = append(finalNodes, node)
		}
	}

	workflow.Nodes = finalNodes

	_, err = w.UpdateWorkflow(workflowId, workflow)

	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateNode updates an existing node
func (w *Workflows) UpdateNode(workflowId string, nodeId string, updateNode N8nNode) (N8nNode, error) {
	if err := w.validateNodeInput(updateNode); err != nil {
		return N8nNode{}, err
	}

	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return N8nNode{}, err
	}

	var finalNodes []N8nNode

	for _, node := range workflow.Nodes {
		if node.Id == nodeId {
			finalNodes = append(finalNodes, w.combineNodes(node, updateNode))
		} else {
			finalNodes = append(finalNodes, node)
		}
	}

	workflow.Nodes = finalNodes

	_, err = w.UpdateWorkflow(workflowId, workflow)

	if err != nil {
		return N8nNode{}, err
	}

	updatedNode, err := w.GetNodeById(workflowId, nodeId)

	if err != nil {
		return N8nNode{}, err
	}

	return updatedNode, nil
}

// combineNodes combines two nodes into one, overwriting the original node with the update node
func (w *Workflows) combineNodes(originalNode N8nNode, updateNode N8nNode) N8nNode {
	if updateNode.Parameters == nil {
		updateNode.Parameters = map[string]interface{}{}
	}

	if updateNode.Credentials == nil {
		updateNode.Credentials = map[string]interface{}{}
	}

	jsonNodeOriginal, err := json.Marshal(originalNode)
	if err != nil {
		return N8nNode{}
	}
	jsonNodeUpdate, err := json.Marshal(updateNode)
	if err != nil {
		return N8nNode{}
	}

	var finalNode N8nNode
	err = json.Unmarshal(jsonNodeOriginal, &finalNode)
	if err != nil {
		return N8nNode{}
	}

	err = json.Unmarshal(jsonNodeUpdate, &finalNode)
	if err != nil {
		return N8nNode{}
	}

	return finalNode
}

// validateNodeInput validates the input for creating or updating a node
func (w *Workflows) validateNodeInput(node N8nNode) error {
	if node.Id != "" {
		return fmt.Errorf("id should be empty when creating or updating a node")
	}

	if node.Type == "" {
		return fmt.Errorf("type should not be empty when creating or updating a node")
	}

	if len(node.Position) != 2 {
		return fmt.Errorf("position should be an array of 2 integers when creating or updating a node")
	}

	if node.Name == "" {
		return fmt.Errorf("name should not be empty when creating or updating a node")
	}

	return nil
}
