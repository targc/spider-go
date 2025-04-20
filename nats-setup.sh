#!/bin/bash

STREAM=local-output
nats stream add --server nats:4222 --user root --password root --storage=memory --replicas=1  --retention=limits --discard=old --max-msgs=-1 --max-msgs-per-subject=-1 --no-allow-rollup --max-bytes=-1 --max-age=1h --max-msg-size=-1 --dupe-window=2m --deny-delete --deny-purge --subjects=$STREAM $STREAM

STREAM=local-input
nats stream add --server nats:4222 --user root --password root --storage=memory --replicas=1  --retention=limits --discard=old --max-msgs=-1 --max-msgs-per-subject=-1 --no-allow-rollup --max-bytes=-1 --max-age=1h --max-msg-size=-1 --dupe-window=2m --deny-delete --deny-purge --subjects=$STREAM $STREAM

# WORKFLOW

CONSUMER_ID=local-workflow
STREAM=local-output
nats consumer add --server nats:4222 --user root --password root --deliver=all --pull --ack=all --replay=instant --max-deliver=-1 --max-pending=0 --no-headers-only --backoff=none --filter= $STREAM $CONSUMER_ID

# WORKER [A]

STREAM=local-input
nats stream add --server nats:4222 --user root --password root --storage=memory --replicas=1  --retention=limits --discard=old --max-msgs=-1 --max-msgs-per-subject=-1 --no-allow-rollup --max-bytes=-1 --max-age=1h --max-msg-size=-1 --dupe-window=2m --deny-delete --deny-purge --subjects=$STREAM $STREAM

CONSUMER_ID=local-worker-test-action-a
STREAM=local-input
nats consumer add --server nats:4222 --user root --password root --deliver=all --pull --ack=all --replay=instant --max-deliver=-1 --max-pending=0 --no-headers-only --backoff=none --filter= $STREAM $CONSUMER_ID

# WORKER [B]

CONSUMER_ID=local-worker-test-action-b
STREAM=local-input
nats consumer add --server nats:4222 --user root --password root --deliver=all --pull --ack=all --replay=instant --max-deliver=-1 --max-pending=0 --no-headers-only --backoff=none --filter= $STREAM $CONSUMER_ID
