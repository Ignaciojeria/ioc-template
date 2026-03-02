# structure

> Project directory structure

```
.
├── .einar.template.json
├── .environment
├── .github
│   └── workflows
│       └── coverage.yml
├── .gitignore
├── .version
├── README.md
├── app
│   ├── adapter
│   │   ├── in
│   │   │   ├── eventbus
│   │   │   │   ├── consumer.go
│   │   │   │   └── consumer_test.go
│   │   │   └── fuegoapi
│   │   │       ├── delete.go
│   │   │       ├── delete_test.go
│   │   │       ├── get.go
│   │   │       ├── get_test.go
│   │   │       ├── patch.go
│   │   │       ├── patch_test.go
│   │   │       ├── post.go
│   │   │       ├── post_test.go
│   │   │       ├── put.go
│   │   │       └── put_test.go
│   │   └── out
│   │       ├── eventbus
│   │       │   ├── publisher.go
│   │       │   └── publisher_test.go
│   │       └── postgres
│   │           ├── postgres_repository.go
│   │           └── postgres_repository_test.go
│   ├── application
│   │   └── usecase
│   │       └── interfaces.go
│   └── shared
│       ├── configuration
│       │   ├── conf.go
│       │   ├── conf_test.go
│       │   ├── parse.go
│       │   └── parse_test.go
│       └── infrastructure
│           ├── eventbus
│           │   ├── factory.go
│           │   ├── gcp_client.go
│           │   ├── gcp_client_test.go
│           │   ├── gcp_publisher.go
│           │   ├── gcp_publisher_test.go
│           │   ├── gcp_subscriber.go
│           │   ├── gcp_subscriber_test.go
│           │   ├── nats_client.go
│           │   ├── nats_client_test.go
│           │   ├── nats_publisher.go
│           │   ├── nats_subscriber.go
│           │   ├── strategy.go
│           │   └── strategy_test.go
│           ├── httpserver
│           │   ├── doc
│           │   │   └── openapi.json
│           │   ├── middleware
│           │   │   ├── request_logger.go
│           │   │   └── request_logger_test.go
│           │   ├── server.go
│           │   └── server_test.go
│           ├── observability
│           │   ├── observability.go
│           │   └── observability_test.go
│           └── postgresql
│               ├── connection.go
│               ├── connection_test.go
│               └── migrations
│                   ├── 000001_initial_schema.down.sql
│                   └── 000001_initial_schema.up.sql
├── cmd
│   └── api
│       └── main.go
├── codecov.yml
├── coverage
├── coverage.out
├── coverage_p.out
├── docker-compose.yml
├── go.mod
├── go.sum
├── scripts
│   └── gen-skills.config.yaml
├── template-generada
│   ├── .environment
│   ├── .gitignore
│   ├── .version
│   ├── app
│   │   ├── adapter
│   │   │   ├── in
│   │   │   │   └── fuegoapi
│   │   │   │       ├── get_account.go
│   │   │   │       ├── get_account_ledger.go
│   │   │   │       ├── get_account_ledger_test.go
│   │   │   │       ├── get_account_test.go
│   │   │   │       ├── get_transaction.go
│   │   │   │       ├── get_transaction_test.go
│   │   │   │       ├── post_account.go
│   │   │   │       ├── post_account_test.go
│   │   │   │       ├── post_transaction.go
│   │   │   │       └── post_transaction_test.go
│   │   │   └── out
│   │   │       └── postgres
│   │   │           ├── account_repository.go
│   │   │           ├── idempotency_store.go
│   │   │           ├── ledger_entry_repository.go
│   │   │           ├── transaction_executor.go
│   │   │           └── transaction_repository.go
│   │   ├── application
│   │   │   └── usecase
│   │   │       ├── create_account.go
│   │   │       ├── create_transaction.go
│   │   │       ├── get_account.go
│   │   │       ├── get_account_ledger.go
│   │   │       └── get_transaction.go
│   │   ├── domain
│   │   │   ├── entity
│   │   │   │   ├── account.go
│   │   │   │   ├── idempotency.go
│   │   │   │   ├── ledger_entry.go
│   │   │   │   └── transaction.go
│   │   │   ├── errors
│   │   │   │   └── errors.go
│   │   │   └── port
│   │   │       ├── account_repository.go
│   │   │       ├── idempotency_store.go
│   │   │       ├── ledger_entry_repository.go
│   │   │       ├── transaction_executor.go
│   │   │       └── transaction_repository.go
│   │   └── shared
│   │       ├── configuration
│   │       │   ├── conf.go
│   │       │   ├── conf_test.go
│   │       │   ├── parse.go
│   │       │   └── parse_test.go
│   │       └── infrastructure
│   │           ├── httpserver
│   │           │   ├── doc
│   │           │   │   └── openapi.json
│   │           │   ├── middleware
│   │           │   │   ├── request_logger.go
│   │           │   │   └── request_logger_test.go
│   │           │   ├── server.go
│   │           │   └── server_test.go
│   │           ├── observability
│   │           │   ├── observability.go
│   │           │   └── observability_test.go
│   │           └── postgresql
│   │               ├── connection.go
│   │               └── migrations
│   │                   ├── 000001_initial_schema.down.sql
│   │                   ├── 000001_initial_schema.up.sql
│   │                   ├── 000002_ledger_schema.down.sql
│   │                   └── 000002_ledger_schema.up.sql
│   ├── cmd
│   │   └── api
│   │       └── main.go
│   ├── docker-compose.yml
│   ├── go.mod
│   ├── go.sum
│   ├── skills-lock.json
│   └── version.go
└── version.go
```
