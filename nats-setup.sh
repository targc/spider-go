#!/bin/bash

# STREAM=local-value-input-test-node-a
# nats stream add --storage=memory --replicas=1  --retention=limits --discard=old --max-msgs=-1 --max-msgs-per-subject=-1 --no-allow-rollup --max-bytes=-1 --max-age=1h --max-msg-size=-1 --dupe-window=2m --deny-delete --deny-purge --subjects=$STREAM $STREAM
#
#
# STREAM=local-value-input-test-node-b
# nats stream add --storage=memory --replicas=1  --retention=limits --discard=old --max-msgs=-1 --max-msgs-per-subject=-1 --no-allow-rollup --max-bytes=-1 --max-age=1h --max-msg-size=-1 --dupe-window=2m --deny-delete --deny-purge --subjects=$STREAM $STREAM

# WORKFLOW

STREAM=local-value-output
nats stream add --server nats:4222 --user root --password root --storage=memory --replicas=1  --retention=limits --discard=old --max-msgs=-1 --max-msgs-per-subject=-1 --no-allow-rollup --max-bytes=-1 --max-age=1h --max-msg-size=-1 --dupe-window=2m --deny-delete --deny-purge --subjects=$STREAM $STREAM

CONSUMER_ID=local-workflow
STREAM=local-value-output
nats consumer add --server nats:4222 --user root --password root --deliver=all --pull --ack=all --replay=instant --max-deliver=-1 --max-pending=0 --no-headers-only --backoff=none --filter= $STREAM $CONSUMER_ID

# WORKER [A]

STREAM=local-value-input
nats stream add --server nats:4222 --user root --password root --storage=memory --replicas=1  --retention=limits --discard=old --max-msgs=-1 --max-msgs-per-subject=-1 --no-allow-rollup --max-bytes=-1 --max-age=1h --max-msg-size=-1 --dupe-window=2m --deny-delete --deny-purge --subjects=$STREAM $STREAM

CONSUMER_ID=local-worker-test-node-a
STREAM=local-value-input
nats consumer add --server nats:4222 --user root --password root --deliver=all --pull --ack=all --replay=instant --max-deliver=-1 --max-pending=0 --no-headers-only --backoff=none --filter= $STREAM $CONSUMER_ID

# WORKER [B]

STREAM=local-value-input
nats stream add --server nats:4222 --user root --password root --storage=memory --replicas=1  --retention=limits --discard=old --max-msgs=-1 --max-msgs-per-subject=-1 --no-allow-rollup --max-bytes=-1 --max-age=1h --max-msg-size=-1 --dupe-window=2m --deny-delete --deny-purge --subjects=$STREAM $STREAM

CONSUMER_ID=local-worker-test-node-b
STREAM=local-value-input
nats consumer add --server nats:4222 --user root --password root --deliver=all --pull --ack=all --replay=instant --max-deliver=-1 --max-pending=0 --no-headers-only --backoff=none --filter= $STREAM $CONSUMER_ID
