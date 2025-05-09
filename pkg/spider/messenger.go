package spider

import (
	"context"
)

type WorkflowMessengerAdapter interface {
	ListenTriggerMessages(ctx context.Context, h func(c TriggerMessageContext, message TriggerMessage) error) error
	ListenOutputMessages(ctx context.Context, h func(c OutputMessageContext, message OutputMessage) error) error
	SendInputMessage(ctx context.Context, message InputMessage) error
	Close(ctx context.Context) error
}

type WorkerMessengerAdapter interface {
	ListenInputMessages(ctx context.Context, h func(c InputMessageContext, message InputMessage) error) error
	SendTriggerMessage(ctx context.Context, message TriggerMessage) error
	SendOutputMessage(ctx context.Context, message OutputMessage) error
	Close(ctx context.Context) error
}
