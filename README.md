# n8n Client for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/kevop-s/n8n-client-go.svg)](https://pkg.go.dev/github.com/kevop-s/n8n-client-go)

`n8n-client-go` is a Go client library for interacting with the [n8n](https://n8n.io/) workflow automation platform's API. This library provides a type-safe interface to manage workflows, nodes, and connections in your n8n instance.

## Features

- Complete workflow management (create, read, update, delete)
- Node management within workflows
- Connection handling between nodes
- Strongly typed data structures for better development experience
- Compatible with the latest n8n versions

## Requirements

- Go 1.24 or higher
- A running n8n instance

## Installation

To install the library, run:

```bash
go get github.com/kevop-s/n8n-client-go
```

## Basic Usage

### Initialize the Client

```go
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

	// Now you can use n8nWorkflows to interact with workflows
}
```

### Usage Examples

#### Get a Workflow by ID

```go
workflow, err := n8nWorkflows.GetWorkflow("workflow-id")
if err != nil {
    log.Fatal("Error getting workflow: ", err)
}
log.Printf("Workflow: %+v", workflow)
```

#### Create a New Workflow

```go
newWorkflow := workflows.N8nWorkflow{
    Name: "My New Workflow",
    Settings: workflows.N8nWorkflowSettings{
        SaveManualExecutions: true,
    },
}

createdWorkflow, err := n8nWorkflows.CreateWorkflow(newWorkflow)
if err != nil {
    log.Fatal("Error creating workflow: ", err)
}
log.Printf("Created workflow: %+v", createdWorkflow)
```

#### Add a Node to a Workflow

```go
newNode := workflows.N8nNode{
    Name: "Example Node",
    Type: "n8n-nodes-base.httpRequest",
    TypeVersion: 1,
    Position: []int{100, 100},
    Parameters: map[string]interface{}{
        "url": "https://api.example.com/data",
        "method": "GET",
    },
}

addedNode, err := n8nWorkflows.AddNode(workflow.Id, newNode)
if err != nil {
    log.Fatal("Error adding node: ", err)
}
log.Printf("Added node: %+v", addedNode)
```

## Project Structure

```
.
├── pkg/
│   ├── client/         # HTTP client and configuration
│   ├── workflows/       # Workflow business logic
│   ├── users/           # User management
│   └── utils/           # Various utilities
└── main.go              # Example implementation
```

## Documentation

For detailed API documentation, see the [GoDoc documentation](https://pkg.go.dev/github.com/kevop-s/n8n-client-go).

## Running Locally

### Prerequisites

1. Go 1.24 or higher

### Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/kevop-s/n8n-client-go.git
   cd n8n-client-go
   ```

2. Set up your n8n instance URL and API key in your environment or in the code.

You can run n8n locally using docker compose:

> Before running the docker compose command, make sure you have a `.env` file in the root directory of the repository. You can copy the `.env.example` file to `.env` and fill in the values.

```bash
docker compose up -d
```

Then add the n8n.example.com domain to your `/etc/hosts` file:

```bash
127.0.0.1 n8n.example.com
```

3. Run the example:
   ```bash
   go run main.go
   ```

## Contributing

Contributions are welcome! Please read the [contribution guidelines](CONTRIBUTING.md) before submitting a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- The [n8n](https://n8n.io/) team for creating an amazing tool
- All contributors who help improve this project
