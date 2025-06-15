package workflows

import (
	"fmt"

	"github.com/kevop-s/n8n-client-go/pkg/utils"
)

type N8nConnection struct {
	SourceNodeName string
	ConnectionType string
	Outputs        []N8nConnectionOutput
}

type N8nConnectionOutput struct {
	OutputIndex               int
	DestinationNodeName       string
	DestinationNodeInputIndex float64
	DestinationNodeInputType  string
}

func (w *Workflows) GetConnectionBySourceNodeName(workflowId string, connectionSourceNodeName string) (N8nConnection, error) {
	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return N8nConnection{}, err
	}

	for _, connection := range workflow.Connections {
		if connection.SourceNodeName == connectionSourceNodeName {
			return connection, nil
		}
	}

	return N8nConnection{}, nil
}

func (w *Workflows) GetConnections(workflowId string) ([]N8nConnection, error) {
	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return nil, err
	}

	return workflow.Connections, nil
}

func (w *Workflows) AddConnection(workflowId string, connection N8nConnection) (N8nConnection, error) {
	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return N8nConnection{}, err
	}

	for _, conn := range workflow.Connections {
		if conn.SourceNodeName == connection.SourceNodeName {
			return N8nConnection{}, fmt.Errorf("connection already exists")
		}
	}

	workflow.Connections = append(workflow.Connections, connection)

	_, err = w.UpdateWorkflow(workflowId, workflow)

	if err != nil {
		return N8nConnection{}, err
	}

	return connection, nil
}

func (w *Workflows) RemoveConnection(workflowId string, connectionSourceNodeName string) (bool, error) {
	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return false, err
	}

	var finalConnections []N8nConnection

	for _, connection := range workflow.Connections {
		if connection.SourceNodeName != connectionSourceNodeName {
			finalConnections = append(finalConnections, connection)
		}
	}

	workflow.Connections = finalConnections

	_, err = w.UpdateWorkflow(workflowId, workflow)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (w *Workflows) UpdateConnection(workflowId string, connectionSourceNodeName string, connection N8nConnection) (N8nConnection, error) {
	workflow, err := w.GetWorkflow(workflowId)

	if err != nil {
		return N8nConnection{}, err
	}

	var finalConnections []N8nConnection

	for _, connection := range workflow.Connections {
		if connection.SourceNodeName != connectionSourceNodeName {
			finalConnections = append(finalConnections, connection)
		}
	}

	workflow.Connections = finalConnections

	_, err = w.UpdateWorkflow(workflowId, workflow)

	if err != nil {
		return N8nConnection{}, err
	}

	return connection, nil
}

func (w *Workflows) ParseConnectionsToObject(workflow map[string]interface{}) ([]N8nConnection, error) {
	var finalConnections []N8nConnection
	var tmpConnection N8nConnection
	var tmpOutput N8nConnectionOutput

	for k, v := range workflow {
		if k == "connections" {
			for sourceNodeName, connection := range v.(map[string]interface{}) {
				tmpConnection.SourceNodeName = sourceNodeName // Run workflow
				for connectionType, outputs := range connection.(map[string]interface{}) {
					tmpConnection.ConnectionType = connectionType // main
					for outputIndex, destinationNodes := range outputs.([]interface{}) {
						tmpOutput.OutputIndex = outputIndex // 0
						for _, destinationNode := range destinationNodes.([]interface{}) {
							tmpOutput.DestinationNodeName = destinationNode.(map[string]interface{})["node"].(string)         // Loop over queries
							tmpOutput.DestinationNodeInputIndex = destinationNode.(map[string]interface{})["index"].(float64) // 0
							tmpOutput.DestinationNodeInputType = destinationNode.(map[string]interface{})["type"].(string)    // main
						}
						tmpConnection.Outputs = append(tmpConnection.Outputs, tmpOutput)
						tmpOutput = N8nConnectionOutput{}
					}
					finalConnections = append(finalConnections, tmpConnection)
					tmpConnection = N8nConnection{}
				}
			}
		}
	}

	return finalConnections, nil
}

func (w *Workflows) ParseConnectionsToMap(connections []N8nConnection) (map[string]interface{}, error) {
	finalConnections := make(map[string]interface{})
	var tmpOutputs = make([]interface{}, 1)
	var tmpOutput = make([]interface{}, 1)

	for _, connection := range connections {
		for _, output := range connection.Outputs {
			if output.DestinationNodeName == "" {
				tmpOutput = append(tmpOutput, make(map[string]interface{}))
			} else {
				tmpOutput = append(tmpOutput, map[string]interface{}{
					"node":  output.DestinationNodeName,
					"index": output.DestinationNodeInputIndex,
					"type":  output.DestinationNodeInputType,
				})
			}

			if len(tmpOutputs) < output.OutputIndex+1 {
				tmpOutputs = append(tmpOutputs, make([]interface{}, output.OutputIndex+1-len(tmpOutputs))...)
			}

			tmpOutputs = append(tmpOutputs[:output.OutputIndex+1], tmpOutputs[output.OutputIndex:]...)
			tmpOutputs[output.OutputIndex] = utils.RemoveEmptyInterfaces(tmpOutput)
			tmpOutput = []interface{}{}
		}
		sourceNodeMap := make(map[string]interface{})
		sourceNodeMap[connection.ConnectionType] = utils.RemoveEmptyInterfaces(tmpOutputs)
		finalConnections[connection.SourceNodeName] = sourceNodeMap
		tmpOutputs = []interface{}{}
	}

	return finalConnections, nil
}
