# Spider Go

Spider Go is a Go-based workflow engine that uses [NATS](https://nats.io/) for messaging and MongoDB for state storage. It provides building blocks for creating workers and workflows that communicate via messages and persist state.

## Project Overview

Spider Go focuses on composing distributed workers into resilient workflows. Workers exchange data through NATS subjects while their progress is persisted to MongoDB so that executions can survive restarts or failures.

## Tech Stack

- **Language:** Go 1.22+
- **Messaging:** NATS JetStream
- **Database:** MongoDB
- **Containerization:** Docker & Docker Compose
- **Optional Integrations:** Slack webhooks, HTTP triggers, cron scheduling

## Project Structure

```
.
├── cmd/            # Entry points for workers, triggers, and workflows
├── examples/       # Sample applications demonstrating Spider Go usage
├── pkg/            # Reusable library code (workflow runtime, adapters, etc.)
├── deploys/        # Deployment manifests and scripts
└── docker-compose.example-basic.yml # Compose file to run the basic example
```

## Getting Started

1. **Install dependencies**
   - [Go](https://go.dev/) 1.22 or later
   - [Docker](https://www.docker.com/) (for running examples)

2. **Run tests**
   ```bash
   go test ./...
   ```

3. **Run the basic example**
   ```bash
   docker compose -f docker-compose.example-basic.yml up --build
   ```
   This starts NATS, MongoDB, and example workers. Set `SLACK_WEBHOOK_URL` in your environment if you want Slack notifications.

4. **Explore more examples**
   See the [examples](examples) directory for additional sample workers and workflows using Spider Go.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

