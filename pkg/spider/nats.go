package spider

import "fmt"

type NatsOutputMessage struct {
	WorkflowNodeID string `json:"workflow_node_id"`
	NodeID         string `json:"node_id"`
	MetaOutput     string `json:"meta_output"`
	Values         Values `json:"values"`
}

type NatsInputMessage struct {
	WorkflowNodeID string `json:"workflow_node_id"`
	NodeID         string `json:"node_id"`
	Values         Values `json:"values"`
}

func buildInputSubject(prefix string) string {
	return fmt.Sprintf("%s-value-input", prefix)
}

func buildOutputSubject(prefix string) string {
	return fmt.Sprintf("%s-value-output", prefix)
}

func buildWorkflowConsumerID(prefix string) string {
	return fmt.Sprintf("%s-workflow", prefix)
}

func buildWorkerConsumerID(prefix, nodeID string) string {
	return fmt.Sprintf("%s-worker-%s", prefix, nodeID)
}
