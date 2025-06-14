package main

import (
	// "fmt"
	"fmt"
	"log"

	"github.com/kevop-s/n8n-client-go/pkg/client"
	// "github.com/kevop-s/n8n-client-go/pkg/users"
	"github.com/kevop-s/n8n-client-go/pkg/workflows"
)

func main() {
	n8nHost := "https://n8n.example.com/api/v1"
	n8nApiKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjZjRmZjNkOC1hMmRlLTRhODAtYmU1Ni04ZjQ5MjM4YjEzMTUiLCJpc3MiOiJuOG4iLCJhdWQiOiJwdWJsaWMtYXBpIiwiaWF0IjoxNzQ5Nzc3OTczfQ.6sllSbs_VFk2EextYD38EwnWi1SnIoLicZIkMQQZZgU"

	n8nClient, err := client.NewClient(&n8nHost, &n8nApiKey)

	if err != nil {
		log.Fatal("Error creating client", err)
	}

	// n8nUsers := users.NewUsers(n8nClient)

	// _, err = n8nUsers.GetUser("kevops@kevops.com")

	// if err != nil {
	// 	log.Fatal("Error getting user ", err)
	// }

	n8nWorkflows := workflows.NewWorkflows(n8nClient)

	workflowDto := workflows.N8nWorkflow{
		Name:       "New workflow (Test)",
		Settings:   workflows.N8nWorkflowSettings{},
		StaticData: "",
	}

	workflow, err := n8nWorkflows.CreateWorkflow(workflowDto)

	if err != nil {
		log.Fatal("Error creating workflow ", err)
	}

	workflowNodeDto := workflows.N8nNode{
		Name:        "When clicking ‘Execute workflow’",
		Type:        "n8n-nodes-base.manualTrigger",
		TypeVersion: 1,
		Position:    []int{0, 0},
		Parameters:  map[string]interface{}{},
		Credentials: map[string]interface{}{},
	}

	workflowNodeDto2 := workflows.N8nNode{
		Name:        "Filter",
		Type:        "n8n-nodes-base.filter",
		TypeVersion: 2.2,
		Position:    []int{200, 0},
		Parameters:  map[string]interface{}{},
		Credentials: map[string]interface{}{},
	}

	_, err = n8nWorkflows.AddNode(workflow.Id, workflowNodeDto)

	if err != nil {
		log.Fatalf("Error adding node %s: %s", workflowNodeDto.Name, err)
	}

	addedNode, err := n8nWorkflows.AddNode(workflow.Id, workflowNodeDto2)

	if err != nil {
		log.Fatalf("Error adding node %s: %s", workflowNodeDto2.Name, err)
	}

	updateNodeDto := workflows.N8nNode{
		Name:             "Filter",
		Type:             "n8n-nodes-base.filter",
		TypeVersion:      2.2,
		Position:         []int{200, 0},
		Parameters:       map[string]interface{}{},
		Credentials:      map[string]interface{}{},
		ExecuteOnce:      true,
		AlwaysOutputData: true,
	}

	_, err = n8nWorkflows.UpdateNode(workflow.Id, addedNode.Id, updateNodeDto)

	if err != nil {
		log.Fatal("Error updating node ", err)
	}

	updatedWorkflowDto := workflows.N8nWorkflow{
		Name: "New workflow (Test Updated)",
		Settings: workflows.N8nWorkflowSettings{
			SaveManualExecutions: true,
		},
	}

	_, err = n8nWorkflows.UpdateWorkflow(workflow.Id, updatedWorkflowDto)

	if err != nil {
		log.Fatal("Error updating workflow ", err)
	}

	fmt.Printf("Workflow updated")

}
