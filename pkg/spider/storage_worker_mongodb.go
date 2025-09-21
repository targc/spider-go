package spider

import (
	"context"

	"github.com/sethvargo/go-envconfig"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func InitMongodDBWorkerStorageAdapter(ctx context.Context) (*MongodDBWorkerStorageAdapter, error) {

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

	a := NewMongodDBWorkerStorageAdapter(client, db)

	return a, nil
}

type MongodDBWorkerStorageAdapter struct {
	client                   *mongo.Client
	workflowActionCollection *mongo.Collection
}

func NewMongodDBWorkerStorageAdapter(client *mongo.Client, db *mongo.Database) *MongodDBWorkerStorageAdapter {
	return &MongodDBWorkerStorageAdapter{
		client:                   client,
		workflowActionCollection: db.Collection("workflow_actions"),
	}
}

func (w *MongodDBWorkerStorageAdapter) GetAllConfigs(ctx context.Context, actionID string) ([]WorkerConfig, error) {

	cur, err := w.workflowActionCollection.Find(
		ctx,
		bson.D{
			{Key: "action_id", Value: actionID},
		},
	)

	if err != nil {
		return nil, err
	}

	var confs []WorkerConfig

	for cur.TryNext(ctx) {

		var workerAction MDWorkflowAction

		err := cur.Decode(&workerAction)

		if err != nil {
			return nil, err
		}

		confs = append(confs, WorkerConfig{
			WorkflowActionID: workerAction.ID,
			TenantID:         workerAction.TenantID,
			WorkflowID:       workerAction.WorkflowID,
			Key:              workerAction.Key,
			Config:           workerAction.Config,
		})
	}

	return confs, nil
}

func (w *MongodDBWorkerStorageAdapter) Close(ctx context.Context) error {
	return w.client.Disconnect(ctx)
}
