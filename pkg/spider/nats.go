package spider

import "fmt"

type NatsOutputMessage struct {
	WorkflowActionID string `json:"workflow_action_id"`
	ActionID         string `json:"action_id"`
	MetaOutput       string `json:"meta_output"`
	Values           Values `json:"values"`
}

type NatsInputMessage struct {
	WorkflowActionID string `json:"workflow_action_id"`
	ActionID         string `json:"action_id"`
	Values           Values `json:"values"`
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

func buildWorkerConsumerID(prefix, actionID string) string {
	return fmt.Sprintf("%s-worker-%s", prefix, actionID)
}
