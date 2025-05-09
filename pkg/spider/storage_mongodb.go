package spider

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sethvargo/go-envconfig"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongodDBWorkflowStorageAdapter struct {
	client                           *mongo.Client
	workflowCollection               *mongo.Collection
	workflowActionCollection         *mongo.Collection
	workflowActionDepCollection      *mongo.Collection
	workflowSessionContextCollection *mongo.Collection
}

type InitMongodDBWorkflowStorageAdapterOpt struct {
	BetaAutoSetupSchema bool
}

func InitMongodDBWorkflowStorageAdapter(ctx context.Context, opt InitMongodDBWorkflowStorageAdapterOpt) (*MongodDBWorkflowStorageAdapter, error) {

	type Env struct {
		MongoDBURI  string `env:"MONGODB_URI,required"`
		MongoDBName string `env:"MONGODB_DB_NAME,required"`
	}

	var env Env

	err := envconfig.Process(ctx, &env)

	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(options.Client().ApplyURI(env.MongoDBURI))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, err
	}

	db := client.Database(env.MongoDBName)

	if opt.BetaAutoSetupSchema {
		err = db.CreateCollection(ctx, "workflows")

		if err != nil {
			return nil, err
		}

		err = db.CreateCollection(ctx, "workflow_actions")

		if err != nil {
			return nil, err
		}

		err = db.CreateCollection(ctx, "workflow_action_deps")

		if err != nil {
			return nil, err
		}

		err = db.CreateCollection(ctx, "workflow_session_contexts")

		if err != nil {
			return nil, err
		}

		_, err = db.Collection("workflow_actions").Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{
				{Key: "key", Value: -1},
				{Key: "workflow_id", Value: -1},
			},
			Options: options.Index().SetUnique(true),
		})

		if err != nil {
			return nil, err
		}

		_, err = db.Collection("workflow_action_deps").Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{
				{Key: "key", Value: -1},
				{Key: "meta_output", Value: -1},
				{Key: "dep_key", Value: -1},
				{Key: "workflow_id", Value: -1},
			},
			Options: options.Index().SetUnique(true),
		})

		if err != nil {
			return nil, err
		}

		_, err = db.Collection("workflow_session_contexts").Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{
				{Key: "session_id", Value: -1},
				{Key: "key", Value: -1},
				{Key: "workflow_id", Value: -1},
			},
			Options: options.Index().SetUnique(true),
		})

		if err != nil {
			return nil, err
		}
	}

	a := NewMongodDBWorkflowStorageAdapter(client, db)

	return a, nil
}

func NewMongodDBWorkflowStorageAdapter(client *mongo.Client, db *mongo.Database) *MongodDBWorkflowStorageAdapter {
	return &MongodDBWorkflowStorageAdapter{
		client:                           client,
		workflowCollection:               db.Collection("workflows"),
		workflowActionCollection:         db.Collection("workflow_actions"),
		workflowActionDepCollection:      db.Collection("workflow_action_deps"),
		workflowSessionContextCollection: db.Collection("workflow_session_contexts"),
	}
}

func (w *MongodDBWorkflowStorageAdapter) AddAction(ctx context.Context, workflowID, key, actionID string, m map[string]Mapper) (*WorkflowAction, error) {

	id, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	wa := MDWorkflowAction{
		ID:         id.String(),
		Key:        key,
		WorkflowID: workflowID,
		ActionID:   actionID,
		Map:        m,
	}

	_, err = w.workflowActionCollection.InsertOne(ctx, wa)

	if err != nil {
		return nil, err
	}

	return &WorkflowAction{
		ID:         wa.ID,
		Key:        wa.Key,
		WorkflowID: wa.WorkflowID,
		ActionID:   wa.ActionID,
		Map:        wa.Map,
	}, nil
}

func (w *MongodDBWorkflowStorageAdapter) AddDep(
	ctx context.Context,
	workflowID,
	key,
	metaOutput,
	depKey string,
) error {
	id, err := uuid.NewV7()

	if err != nil {
		return err
	}

	dep := MDWorkflowActionDep{
		ID:         id.String(),
		WorkflowID: workflowID,
		Key:        key,
		MetaOutput: metaOutput,
		DepKey:     depKey,
	}

	_, err = w.workflowActionDepCollection.InsertOne(ctx, dep)

	if err != nil {
		return err
	}

	return nil
}

func (w *MongodDBWorkflowStorageAdapter) QueryWorkflowAction(ctx context.Context, workflowID, key string) (*WorkflowAction, error) {

	result := w.workflowActionCollection.FindOne(
		ctx,
		bson.D{
			{Key: "workflow_id", Value: workflowID},
			{Key: "key", Value: key},
		},
	)

	err := result.Err()

	if err != nil {
		return nil, err
	}

	var wa MDWorkflowAction

	err = result.Decode(&wa)

	if err != nil {
		return nil, err
	}

	return &WorkflowAction{
		ID:         wa.ID,
		Key:        wa.Key,
		WorkflowID: wa.WorkflowID,
		ActionID:   wa.ActionID,
		Map:        wa.Map,
	}, nil
}

func (w *MongodDBWorkflowStorageAdapter) QueryWorkflowActionDependencies(ctx context.Context, workflowID, key, metaOutput string) ([]WorkflowAction, error) {

	cur, err := w.workflowActionDepCollection.Find(
		ctx,
		bson.D{
			{Key: "workflow_id", Value: workflowID},
			{Key: "key", Value: key},
			{Key: "meta_output", Value: metaOutput},
		},
	)

	if err != nil {
		return nil, err
	}

	var deps []MDWorkflowActionDep

	for cur.TryNext(ctx) {

		var dep MDWorkflowActionDep

		err := cur.Decode(&dep)

		if err != nil {
			return nil, err
		}

		deps = append(deps, dep)
	}

	var depActions []WorkflowAction

	for _, dep := range deps {
		depAction, err := w.QueryWorkflowAction(ctx, workflowID, dep.DepKey)

		if err != nil {
			continue
		}

		depActions = append(depActions, *depAction)
	}

	return depActions, nil
}

func (w *MongodDBWorkflowStorageAdapter) TryAddSessionContext(ctx context.Context, workflowID, sessionID, key string, value map[string]interface{}) (map[string]map[string]interface{}, error) {

	id, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	newSess := MDWorkflowSessionContext{
		ID:         id.String(),
		WorkflowID: workflowID,
		SessionID:  sessionID,
		Key:        key,
		Value:      value,
	}

	_, err = w.workflowSessionContextCollection.InsertOne(ctx, newSess)

	if err != nil {
		return nil, err
	}

	cur, err := w.workflowSessionContextCollection.Find(
		ctx,
		bson.D{
			{Key: "workflow_id", Value: workflowID},
			{Key: "session_id", Value: sessionID},
		},
	)

	if err != nil {
		return nil, err
	}

	var sessCtxs []MDWorkflowSessionContext

	for cur.TryNext(ctx) {

		var sessCtx MDWorkflowSessionContext

		err := cur.Decode(&sessCtx)

		if err != nil {
			return nil, err
		}

		sessCtxs = append(sessCtxs, sessCtx)
	}

	results := map[string]map[string]interface{}{}

	for _, sessCtx := range sessCtxs {
		b, err := json.Marshal(sessCtx.Value)

		if err != nil {
			return nil, err
		}

		result := map[string]interface{}{}

		err = json.Unmarshal(b, &result)

		if err != nil {
			return nil, err
		}

		results[sessCtx.Key] = result
	}

	return results, nil
}

func (w *MongodDBWorkflowStorageAdapter) Close(ctx context.Context) error {
	return w.client.Disconnect(ctx)
}

type MDWorkflowAction struct {
	ID         string            `bson:"_id"`
	Key        string            `bson:"key"`         // Composite unique index
	WorkflowID string            `bson:"workflow_id"` // Composite unique index
	ActionID   string            `bson:"action_id"`
	Map        map[string]Mapper `bson:"map"`
}

type MDWorkflowActionDep struct {
	ID         string `bson:"_id"`
	WorkflowID string `bson:"workflow_id"` // Composite unique index
	Key        string `bson:"key"`         // Composite unique index
	MetaOutput string `bson:"meta_output"` // Composite unique index
	DepKey     string `bson:"dep_key"`     // Composite unique index
}

type MDWorkflowSessionContext struct {
	ID         string                 `bson:"_id"`
	WorkflowID string                 `bson:"workflow_id"` // Composite unique index
	SessionID  string                 `bson:"session_id"`  // Composite unique index
	Key        string                 `bson:"key"`         // Composite unique index
	Value      map[string]interface{} `bson:"value"`
}
