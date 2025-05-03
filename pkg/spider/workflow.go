package spider

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/expr-lang/expr"
	"golang.org/x/sync/errgroup"
)

type Workflow struct {
	messenger WorkflowMessengerAdapter
	storage   WorkflowStorageAdapter
}

func InitWorkflow(
	messenger WorkflowMessengerAdapter,
	storage WorkflowStorageAdapter,
) *Workflow {
	return &Workflow{
		messenger,
		storage,
	}
}

func InitDefaultWorkflow(
	ctx context.Context,
	storage WorkflowStorageAdapter,
) (*Workflow, error) {
	messenger, err := InitNATSWorkflowMessengerAdapter(ctx)

	if err != nil {
		return nil, err
	}

	return &Workflow{
		messenger,
		storage,
	}, nil
}

func (w *Workflow) Run(ctx context.Context) error {

	err := w.messenger.ListenOutputMessages(ctx, func(c OutputMessageContext, m OutputMessage) error {

		workflowAction, err := w.storage.QueryWorkflowAction(c.Context, m.WorkflowID, m.Key)

		if err != nil {
			slog.Error(
				"QueryWorkflowAction failed",
				slog.Any("error", err.Error()),
				slog.Any("workflow_id", m.WorkflowID),
				slog.Any("key", m.Key),
			)

			return err
		}

		_ = workflowAction // TODO:

		wvalues := map[string]interface{}{}

		err = json.Unmarshal([]byte(m.Values), &wvalues)

		if err != nil {
			slog.Error("unmarshal value failed", slog.Any("error", err.Error()))
			return err
		}

		newContextVal := map[string]interface{}{
			"output": wvalues,
		}

		wcontext, err := w.storage.TryAddSessionContext(ctx, m.SessionID, m.Key, newContextVal)

		if err != nil {
			slog.Error("TryAddSessionContext failed", slog.Any("error", err.Error()))
			return err
		}

		deps, err := w.storage.QueryWorkflowActionDependencies(c.Context, m.WorkflowID, m.Key, m.MetaOutput)

		if err != nil {
			slog.Error("QueryWorkflowActionDependencies failed", slog.Any("error", err.Error()))
			return err
		}

		eg := errgroup.Group{}

		eg.SetLimit(10)

		for _, dep := range deps {
			eg.Go(func() error {

				mapper, err := w.storage.QueryWorkflowActionMapper(ctx, m.WorkflowID, m.Key, m.MetaOutput, dep.Key)

				if err != nil {
					slog.Error("QueryWorkflowActionMapper failed", slog.Any("error", err.Error()))
					return err
				}

				nextInput, err := ex(wcontext, mapper)

				if err != nil {
					slog.Error("ex failed", slog.Any("error", err.Error()))
					return err
				}

				nextInputb, err := json.Marshal(nextInput)

				if err != nil {
					slog.Error("marshal next input failed", slog.Any("error", err.Error()))
					return err
				}

				err = w.messenger.SendInputMessage(ctx, InputMessage{
					SessionID:  m.SessionID,
					WorkflowID: dep.WorkflowID,
					// TODO
					// WorkflowActionID: dep.ID,
					Key:      dep.Key,
					ActionID: dep.ActionID,
					Values:   string(nextInputb),
				})

				if err != nil {
					slog.Error("sent input message failed", slog.Any("error", err.Error()))
					return err
				}

				return nil
			})
		}

		err = eg.Wait()

		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (w *Workflow) Close(ctx context.Context) error {

	err := w.messenger.Close(ctx)

	if err != nil {
		return err
	}

	err = w.storage.Close(ctx)

	if err != nil {
		return err
	}

	return nil
}

func ex(env map[string]map[string]interface{}, mapping map[string]Mapper) (map[string]interface{}, error) {
	output := map[string]interface{}{}

	for k, v := range mapping {

		if len(v.Value) == 0 {
			output[k] = ""
			continue
		}

		if v.Mode == MapperModeFixed {
			output[k] = v.Value
			continue
		}

		expression := v.Value

		slog.Info("executing expression", slog.String("expression", expression))

		program, err := expr.Compile(expression, expr.Env(env))

		if err != nil {
			// output[k] = fmt.Sprintf("<compile error: %v>", err)
			// continue
			return nil, fmt.Errorf("error on expression %v: %s", expression, err.Error())
		}

		slog.Info("executing program", slog.String("disassemble", program.Disassemble()))

		result, err := expr.Run(program, env)

		if err != nil {
			// output[k] = fmt.Sprintf("<runtime error: %v>", err)
			// continue
			return nil, fmt.Errorf("error on expression %v: %s", expression, err.Error())
		}

		output[k] = result
	}

	return output, nil
}
