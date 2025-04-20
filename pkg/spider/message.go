package spider

import (
	"context"
	"time"
)

type InputMessageContext struct {
	Context    context.Context
	Timestamp  time.Time
	SendOutput func(metaOutput string, values string) error
}

type InputMessage struct {
	SessionID        string
	WorkflowActionID string
	Values           string
}

type OutputMessageContext struct {
	Context   context.Context
	Timestamp time.Time
}

type OutputMessage struct {
	SessionID        string
	WorkflowActionID string
	MetaOutput       string
	Values           string
}
