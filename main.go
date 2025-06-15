package main

import (
	"log"

	"github.com/kevop-s/n8n-client-go/pkg/client"
	"github.com/kevop-s/n8n-client-go/pkg/workflows"
)

func main() {
	n8nHost := "https://your-n8n-instance.com/api/v1"
	n8nApiKey := "your-api-key"

	n8nClient, err := client.NewClient(&n8nHost, &n8nApiKey)
	if err != nil {
		log.Fatal("Error creating client: ", err)
	}

	// Initialize workflows handler
	n8nWorkflows := workflows.NewWorkflows(n8nClient)

	workflow, err := n8nWorkflows.GetWorkflow("workflow-id")
	if err != nil {
		log.Fatal("Error getting workflow: ", err)
	}
	log.Printf("Workflow: %+v", workflow)

	// Now you can use n8nWorkflows to interact with workflows
}
