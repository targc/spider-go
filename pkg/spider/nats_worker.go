package spider

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sethvargo/go-envconfig"
	"github.com/targc/xnats-go"
	"golang.org/x/sync/errgroup"
)

type NATSWorkerMessengerAdapter struct {
	nc               *xnats.XNats
	p                *xnats.Producer
	c                *xnats.Consumer
	cctx             jetstream.ConsumeContext
	natsStreamPrefix string
	nodeID           string
}

var _ WorkerMessengerAdapter = &NATSWorkerMessengerAdapter{}

func InitNATSWorkerMessengerAdapter(ctx context.Context, nodeID string) (*NATSWorkerMessengerAdapter, error) {
	type Env struct {
		NATSHost             string `env:"NATS_HOST,required"`
		NATSPort             int    `env:"NATS_PORT,required"`
		NATSUser             string `env:"NATS_USER,required"`
		NATSPassword         string `env:"NATS_PASSWORD,required"`
		NATSStreamPrefix     string `env:"NATS_STREAM_PREFIX,required"`
		NATSConsumerIDPrefix string `env:"NATS_CONSUMER_ID_PREFIX,required"`
	}

	var env Env

	err := envconfig.Process(ctx, &env)

	if err != nil {
		return nil, err
	}

	nc, err := xnats.Connect(xnats.ConnectOpt{
		Host:     env.NATSHost,
		Port:     env.NATSPort,
		User:     env.NATSUser,
		Password: env.NATSPassword,
	})

	if err != nil {
		return nil, err
	}

	p := nc.Producer()

	inputStream := buildInputSubject(env.NATSStreamPrefix)
	consumerID := buildWorkerConsumerID(env.NATSConsumerIDPrefix, nodeID)

	slog.Info(
		"worker",
		slog.String("input_stream", inputStream),
		slog.String("consumer_id", consumerID),
		slog.String("node_id", nodeID),
	)

	c, err := nc.Consumer(ctx, inputStream, consumerID)

	if err != nil {
		return nil, err
	}

	adapter := NATSWorkerMessengerAdapter{
		nc:               nc,
		p:                p,
		c:                c,
		natsStreamPrefix: env.NATSStreamPrefix,
		nodeID:           nodeID,
	}

	return &adapter, nil
}

func (m *NATSWorkerMessengerAdapter) ListenInputMessages(ctx context.Context, h func(c InputMessageContext, message InputMessage) error) error {

	if m.cctx != nil {
		return errors.New("cannot re-initialize")
	}

	ictx := context.Background()

	eg := errgroup.Group{}

	eg.SetLimit(50)

	cctx, err := m.c.Consume(func(msg jetstream.Msg) {
		eg.Go(func() error {
			msg.Ack()

			slog.Info(
				"received input",
				slog.String("b", string(msg.Data())),
			)

			metadata, err := msg.Metadata()

			if err != nil {
				// TODO:
				return err
			}

			var b NatsOutputMessage

			err = json.Unmarshal(msg.Data(), &b)

			if err != nil {
				// TODO:
				return err
			}

			if b.NodeID != m.nodeID {
				return nil
			}

			err = h(
				InputMessageContext{
					Context:   ictx,
					Timestamp: metadata.Timestamp,
				},
				InputMessage{
					WorkflowNodeID: b.WorkflowNodeID,
					NodeID:         b.NodeID,
					Values:         b.Values,
				},
			)

			if err != nil {
				// TODO:
				return err
			}

			return nil
		})
	})

	if err != nil {
		return err
	}

	m.cctx = cctx

	<-ctx.Done()

	return nil
}

func (m *NATSWorkerMessengerAdapter) SendOutputMessage(ctx context.Context, message OutputMessage) error {
	subject := buildOutputSubject(m.natsStreamPrefix)

	b, err := json.Marshal(NatsOutputMessage{
		WorkflowNodeID: message.WorkflowNodeID,
		NodeID:         message.NodeID,
		MetaOutput:     message.MetaOutput,
		Values:         message.Values,
	})

	if err != nil {
		return err
	}

	slog.Info(
		"sent output",
		slog.String("subject", subject),
		slog.String("b", string(b)),
	)

	err = m.p.Produce(ctx, subject, b)

	if err != nil {
		return err
	}

	return nil
}

func (m *NATSWorkerMessengerAdapter) Close(ctx context.Context) error {
	m.cctx.Stop()
	m.nc.Close()
	return nil
}
