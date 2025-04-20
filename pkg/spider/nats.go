package spider

import "fmt"

type NatsOutputMessage struct {
	SessionID        string `json:"session_id"`
	WorkflowActionID string `json:"workflow_action_id"`
	ActionID         string `json:"action_id"`
	MetaOutput       string `json:"meta_output"`
	Values           string `json:"values"`
}

type NatsInputMessage struct {
	SessionID        string `json:"session_id"`
	WorkflowActionID string `json:"workflow_action_id"`
	ActionID         string `json:"action_id"`
	Values           string `json:"values"`
}

func buildInputSubject(prefix string) string {
	return fmt.Sprintf("%s-input", prefix)
}

func buildOutputSubject(prefix string) string {
	return fmt.Sprintf("%s-output", prefix)
}

func buildWorkflowConsumerID(prefix string) string {
	return fmt.Sprintf("%s-workflow", prefix)
}

func buildWorkerConsumerID(prefix, actionID string) string {
	return fmt.Sprintf("%s-worker-%s", prefix, actionID)
}
