package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"spider-go/pkg/spider"

	"github.com/google/uuid"
)

type MockWorkflowStorageAdapter struct {
	// workflow id -> key -> action
	actions map[string]map[string]spider.WorkflowAction

	// workflow id -> key -> meta output -> dep key -> input field -> mapper
	depsmap map[string]map[string]map[string]map[string]map[string]spider.Mapper

	// workflow id -> session id -> key -> action
	sessions map[string]map[string]map[string]interface{}
}

func NewMockWorkflowStorageAdapter() *MockWorkflowStorageAdapter {
	return &MockWorkflowStorageAdapter{
		actions:  map[string]map[string]spider.WorkflowAction{},
		depsmap:  map[string]map[string]map[string]map[string]map[string]spider.Mapper{},
		sessions: map[string]map[string]map[string]interface{}{},
	}
}

func (w *MockWorkflowStorageAdapter) AddAction(ctx context.Context, workflowID, key, actionID string) (*spider.WorkflowAction, error) {

	_, ok := w.actions[workflowID]

	if !ok {
		w.actions[workflowID] = map[string]spider.WorkflowAction{}
	}

	_, ok = w.actions[workflowID][key]

	if ok {
		return nil, errors.New("duplicated key")
	}

	id, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	workflowAction := spider.WorkflowAction{
		ID:         id.String(),
		Key:        key,
		WorkflowID: workflowID,
		ActionID:   actionID,
	}

	w.actions[workflowID][key] = workflowAction

	log.Printf("workflow action: %#v\n", workflowAction)

	return &workflowAction, nil
}

func (w *MockWorkflowStorageAdapter) AddDep(
	ctx context.Context,
	workflowID,
	key,
	metaOutput,
	key2 string,
	m map[string]spider.Mapper,
) (*spider.WorkflowAction, error) {

	_, ok := w.depsmap[workflowID]

	if !ok {
		w.depsmap[workflowID] = map[string]map[string]map[string]map[string]spider.Mapper{}
	}

	_, ok = w.depsmap[workflowID][key]

	if !ok {
		w.depsmap[workflowID][key] = map[string]map[string]map[string]spider.Mapper{}
	}

	_, ok = w.depsmap[workflowID][key][metaOutput]

	if !ok {
		w.depsmap[workflowID][key][metaOutput] = map[string]map[string]spider.Mapper{}
	}

	w.depsmap[workflowID][key][metaOutput][key2] = m

	return nil, nil
}

func (w *MockWorkflowStorageAdapter) QueryWorkflowAction(ctx context.Context, workflowID, key string) (*spider.WorkflowAction, error) {

	_, ok := w.actions[workflowID]

	if !ok {
		return nil, errors.New("not found workflow")
	}

	wa, ok := w.actions[workflowID][key]

	if !ok {
		return nil, errors.New("not found workflow action")
	}

	return &wa, nil
}

func (w *MockWorkflowStorageAdapter) QueryWorkflowActionDependencies(ctx context.Context, workflowID, key, metaOutput string) ([]spider.WorkflowAction, error) {

	_, ok := w.depsmap[workflowID]

	if !ok {
		return nil, errors.New("not found workflow")
	}

	_, ok = w.depsmap[workflowID][key]

	if !ok {
		return nil, errors.New("not found workflow action")
	}

	deps, ok := w.depsmap[workflowID][key][metaOutput]

	if !ok {
		return nil, errors.New("not found workflow action by meta output")
	}

	var depActions []spider.WorkflowAction

	for key := range deps {
		depAction, err := w.QueryWorkflowAction(ctx, workflowID, key)

		if err != nil {
			continue
		}

		depActions = append(depActions, *depAction)
	}

	return depActions, nil
}

func (w *MockWorkflowStorageAdapter) QueryWorkflowActionMapper(ctx context.Context, workflowID, key, metaOutput, key2 string) (map[string]spider.Mapper, error) {

	_, ok := w.depsmap[workflowID]

	if !ok {
		return nil, errors.New("not found workflow")
	}

	_, ok = w.depsmap[workflowID][key]

	if !ok {
		return nil, errors.New("not found workflow action")
	}

	_, ok = w.depsmap[workflowID][key][metaOutput]

	if !ok {
		return nil, errors.New("not found workflow action by meta output")
	}

	m, ok := w.depsmap[workflowID][key][metaOutput][key2]

	if !ok {
		return nil, errors.New("not found mapper")
	}

	return m, nil
}

func (w *MockWorkflowStorageAdapter) TryAddSessionContext(ctx context.Context, sessionID, contextKey string, contextValue map[string]interface{}) (map[string]map[string]interface{}, error) {

	_, ok := w.sessions[sessionID]

	if !ok {
		w.sessions[sessionID] = map[string]map[string]interface{}{}
	}

	_, ok = w.sessions[sessionID][contextKey]

	if ok {
		return nil, errors.New("context key already added")
	}

	w.sessions[sessionID][contextKey] = contextValue

	return w.sessions[sessionID], nil
}

func (w *MockWorkflowStorageAdapter) Close(ctx context.Context) error {
	w.actions = map[string]map[string]spider.WorkflowAction{}
	w.depsmap = map[string]map[string]map[string]map[string]map[string]spider.Mapper{}
	w.sessions = map[string]map[string]map[string]interface{}{}

	return nil
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	mockStorage := NewMockWorkflowStorageAdapter()

	workflowID := "wa"

	_, err := mockStorage.AddAction(ctx, workflowID, "a1", "test-action-a")

	if err != nil {
		panic(err)
	}

	_, err = mockStorage.AddAction(ctx, workflowID, "a2", "test-action-b")

	if err != nil {
		panic(err)
	}

	_, err = mockStorage.AddAction(ctx, workflowID, "a3", "test-action-b")

	if err != nil {
		panic(err)
	}

	_, err = mockStorage.AddDep(ctx, workflowID, "a1", "triggered", "a2", map[string]spider.Mapper{
		"value": {
			Mode:  spider.MapperModeExpression,
			Value: "a1.output.value + '_updatedx1'",
		},
	})

	if err != nil {
		panic(err)
	}

	_, err = mockStorage.AddDep(ctx, workflowID, "a2", "success", "a3", map[string]spider.Mapper{
		"value": {
			Mode:  spider.MapperModeExpression,
			Value: "a2.output.value + '_updatedx2'",
		},
	})

	if err != nil {
		panic(err)
	}

	worflow, err := spider.InitDefaultWorkflow(ctx, mockStorage)

	if err != nil {
		panic(err)
	}

	go worflow.Run(ctx)

	nctx, ncancel := signal.NotifyContext(ctx, os.Interrupt)
	defer ncancel()

	<-nctx.Done()

	cancel()
	_ = worflow.Close(ctx)
}
