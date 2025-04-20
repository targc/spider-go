package spider

import (
	"context"
	"time"
)

type OutputMessageContext struct {
	Context   context.Context
	Timestamp time.Time
}

type OutputMessage struct {
	WorkflowActionID string
	ActionID         string
	MetaOutput       string
	Values           Values
}

type InputMessageContext struct {
	Context   context.Context
	Timestamp time.Time
}

type InputMessage struct {
	WorkflowActionID string
	ActionID         string
	Values           Values
}
