package spider

import "fmt"

type NatsOutputMessage struct {
	SessionID  string `json:"session_id"`
	WorkflowID string `json:"workflow_id"`
	// TODO
	// WorkflowActionID string `json:"workflow_action_id"`
	MetaOutput string `json:"meta_output"`
	Key        string `json:"key"`
	ActionID   string `json:"action_id"`
	Values     string `json:"values"`
}

func (n NatsOutputMessage) FromOutputMessage(message OutputMessage) NatsOutputMessage {
	return NatsOutputMessage{
		SessionID:  message.SessionID,
		WorkflowID: message.WorkflowID,
		// TODO:
		// WorkflowActionID: message.WorkflowActionID,
		MetaOutput: message.MetaOutput,
		Key:        message.Key,
		ActionID:   message.ActionID,
		Values:     message.Values,
	}
}

func (n *NatsOutputMessage) ToOutputMessage() OutputMessage {
	return OutputMessage{
		SessionID:  n.SessionID,
		WorkflowID: n.WorkflowID,
		// TODO
		// WorkflowActionID: b.WorkflowActionID,
		MetaOutput: n.MetaOutput,
		Key:        n.Key,
		ActionID:   n.ActionID,
		Values:     n.Values,
	}
}

type NatsInputMessage struct {
	SessionID  string `json:"session_id"`
	WorkflowID string `json:"workflow_id"`
	// TODO
	// WorkflowActionID string `json:"workflow_action_id"`
	Key      string `json:"key"`
	ActionID string `json:"action_id"`
	Values   string `json:"values"`
}

func (n *NatsInputMessage) ToInputMessage() InputMessage {
	return InputMessage{
		SessionID:  n.SessionID,
		WorkflowID: n.WorkflowID,
		// TODO
		// WorkflowActionID: n.WorkflowActionID,
		Key:      n.Key,
		ActionID: n.ActionID,
		Values:   n.Values,
	}
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
