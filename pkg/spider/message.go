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
	WorkflowNodeID string
	NodeID         string
	MetaOutput     string
	Values         Values
}

type InputMessageContext struct {
	Context   context.Context
	Timestamp time.Time
}

type InputMessage struct {
	WorkflowNodeID string
	NodeID         string
	Values         Values
}
