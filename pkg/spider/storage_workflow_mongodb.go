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
			// return nil, err
		}

		err = db.CreateCollection(ctx, "workflow_actions")

		if err != nil {
			// return nil, err
		}

		err = db.CreateCollection(ctx, "workflow_action_deps")

		if err != nil {
			// return nil, err
		}

		err = db.CreateCollection(ctx, "workflow_session_contexts")

		if err != nil {
			// return nil, err
		}

		_, err = db.Collection("workflow_actions").Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{
				{Key: "key", Value: -1},
				{Key: "tenant_id", Value: -1},
				{Key: "workflow_id", Value: -1},
			},
			Options: options.Index().SetUnique(true),
		})

		if err != nil {
			// return nil, err
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
			// return nil, err
		}

		_, err = db.Collection("workflow_session_contexts").Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{
				{Key: "session_id", Value: -1},
				{Key: "task_id", Value: -1},
				{Key: "workflow_id", Value: -1},
			},
			Options: options.Index().SetUnique(true),
		})

		if err != nil {
			// return nil, err
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

func (w *MongodDBWorkflowStorageAdapter) AddAction(ctx context.Context, tenantID, workflowID, key, actionID string, conf map[string]string, m map[string]Mapper, meta map[string]string) (*WorkflowAction, error) {

	id, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	wa := MDWorkflowAction{
		ID:         id.String(),
		Key:        key,
		TenantID:   tenantID,
		WorkflowID: workflowID,
		ActionID:   actionID,
		Config:     conf,
		Map:        m,
		Meta:       meta,
		Disabled:   false,
	}

	_, err = w.workflowActionCollection.InsertOne(ctx, wa)

	if err != nil {
		return nil, err
	}

	return &WorkflowAction{
		ID:         wa.ID,
		Key:        wa.Key,
		TenantID:   wa.TenantID,
		WorkflowID: wa.WorkflowID,
		ActionID:   wa.ActionID,
		Map:        wa.Map,
		Meta:       wa.Meta,
		Disabled:   wa.Disabled,
	}, nil
}

func (w *MongodDBWorkflowStorageAdapter) AddDep(
	ctx context.Context,
	tenantID,
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

func (w *MongodDBWorkflowStorageAdapter) QueryWorkflowAction(ctx context.Context, tenantID, workflowID, key string) (*WorkflowAction, error) {

	result := w.workflowActionCollection.FindOne(
		ctx,
		bson.D{
			{Key: "tenant_id", Value: tenantID},
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
		TenantID:   wa.TenantID,
		WorkflowID: wa.WorkflowID,
		ActionID:   wa.ActionID,
		Map:        wa.Map,
		Meta:       wa.Meta,
		Disabled:   wa.Disabled,
	}, nil
}

func (w *MongodDBWorkflowStorageAdapter) QueryWorkflowActionDependencies(ctx context.Context, tenantID, workflowID, key, metaOutput string) ([]WorkflowAction, error) {

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
		depAction, err := w.QueryWorkflowAction(ctx, tenantID, workflowID, dep.DepKey)

		if err != nil {
			continue
		}

		depActions = append(depActions, *depAction)
	}

	return depActions, nil
}

func (w *MongodDBWorkflowStorageAdapter) GetSessionContext(ctx context.Context, workflowID, sessionID, taskID string) (map[string]map[string]interface{}, error) {
	result := w.workflowSessionContextCollection.FindOne(
		ctx,
		bson.D{
			{Key: "workflow_id", Value: workflowID},
			{Key: "session_id", Value: sessionID},
			{Key: "task_id", Value: taskID},
		},
	)

	err := result.Err()

	if err != nil {
		return nil, err
	}

	var sessCtx MDWorkflowSessionContext

	err = result.Decode(&sessCtx)

	if err != nil {
		return nil, err
	}

	valb, err := json.Marshal(sessCtx.Value)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(valb, &sessCtx.Value)

	if err != nil {
		return nil, err
	}

	return sessCtx.Value, nil
}

func (w *MongodDBWorkflowStorageAdapter) CreateSessionContext(ctx context.Context, workflowID, sessionID, taskID string, value map[string]map[string]interface{}) error {
	id, err := uuid.NewV7()

	if err != nil {
		return err
	}

	newSess := MDWorkflowSessionContext{
		ID:         id.String(),
		WorkflowID: workflowID,
		SessionID:  sessionID,
		TaskID:     taskID,
		Value:      value,
	}

	_, err = w.workflowSessionContextCollection.InsertOne(ctx, newSess)

	if err != nil {
		return err
	}

	return nil
}

func (w *MongodDBWorkflowStorageAdapter) DeleteSessionContext(ctx context.Context, workflowID, sessionID, taskID string) error {
	// _, err := w.workflowSessionContextCollection.DeleteOne(
	// 	ctx,
	// 	bson.D{
	// 		{Key: "workflow_id", Value: workflowID},
	// 		{Key: "session_id", Value: sessionID},
	// 		{Key: "task_id", Value: taskID},
	// 	},
	// )
	//
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (w *MongodDBWorkflowStorageAdapter) DisableWorkflowAction(ctx context.Context, tenantID, workflowID, key string) error {

	_, err := w.workflowActionCollection.UpdateOne(
		ctx,
		bson.D{
			{Key: "tenant_id", Value: tenantID},
			{Key: "workflow_id", Value: workflowID},
			{Key: "key", Value: key},
		},
		bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "disabled", Value: true},
				},
			},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (w *MongodDBWorkflowStorageAdapter) Close(ctx context.Context) error {
	return w.client.Disconnect(ctx)
}

type MDWorkflowAction struct {
	ID         string            `bson:"_id"`
	Key        string            `bson:"key"`         // Composite unique index
	TenantID   string            `bson:"tenant_id"`   // Composite unique index
	WorkflowID string            `bson:"workflow_id"` // Composite unique index
	ActionID   string            `bson:"action_id"`
	Config     map[string]string `bson:"config"`
	Map        map[string]Mapper `bson:"map"`
	Meta       map[string]string `bson:"meta,omitempty"`
	Disabled   bool              `bson:"disabled"`
}

type MDWorkflowActionDep struct {
	ID         string `bson:"_id"`
	WorkflowID string `bson:"workflow_id"` // Composite unique index
	Key        string `bson:"key"`         // Composite unique index
	MetaOutput string `bson:"meta_output"` // Composite unique index
	DepKey     string `bson:"dep_key"`     // Composite unique index
}

type MDWorkflowSessionContext struct {
	ID         string                            `bson:"_id"`
	WorkflowID string                            `bson:"workflow_id"` // Composite unique index
	SessionID  string                            `bson:"session_id"`  // Composite unique index
	TaskID     string                            `bson:"task_id"`     // Composite unique index
	Value      map[string]map[string]interface{} `bson:"value"`
}
