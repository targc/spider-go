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
	MetaOutput       string
	Values           string
}

type InputMessageContext struct {
	Context    context.Context
	Timestamp  time.Time
	SendOutput func(metaOutput string, values string) error
}

type InputMessage struct {
	WorkflowActionID string
	Values           string
}
