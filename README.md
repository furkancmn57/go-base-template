# go-base-template

Horizontal-layer Go API template (NotificationApi-style layout). Single binary,
no repository, no outbox, no `cmd/`/`internal/`.

Rules: [`.cursor/rules/go-artictech.mdc`](.cursor/rules/go-artictech.mdc) ·
Naming: [`.cursor/STRUCTURE.md`](.cursor/STRUCTURE.md).

## Stack

Fiber v2 · GORM/Postgres · Redis · RabbitMQ · ozzo-validation · OpenAPI (swaggo) · optional GraphQL

## Layout

```text
src/main.go
src/config/  extensions/  common/  constants/  enums/  interfaces/
src/controllers/v1/  graphql/   # GraphQL optional; same services as REST
src/services/{resource}/ (+ validations/)  services/cache/
src/data/{entities,mappings}/  data/postgres.go  data/migrate.go
src/models/{requests,responses}/
src/events/{domain}/  consumers/{domain}/
src/extensions/  interfaces/
```

## Run

```bash
cp src/config/env/.env.example .env
docker compose up -d postgres redis rabbitmq
make run
```

| URL | |
|-----|--|
| API | http://localhost:8080/api/v1 |
| Health | http://localhost:8080/health |
| OpenAPI | http://localhost:8080/openapi/index.html |
| GraphQL | http://localhost:8080/graphql (set `GRAPHQL_ENABLED=true`) |

## Make

`run` · `build` · `openapi` · `docker-up` · `docker-down` · `test` · `vet`

After changing controller annotations, run `make openapi` (needs `go install github.com/swaggo/swag/cmd/swag@latest`).
