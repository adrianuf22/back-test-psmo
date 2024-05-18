# back-test-psmo

A Golang REST API application using `http.ServeMux` router, `pgx` Postgres driver, `godotenv` with `envconfig` handle configuration.

Requirements:
- Go 1.22+ (For build only)
- Docker 24+ (Compose 2.27+ preferable)

---

## Run

```shell
docker-compose up
```

> It should works: curl -i --request GET --url http://0.0.0.0:3000/v1/health

> Note: `docker-compose build` should be runned every a new change had been made. The current version of development environment of this application hasn't any kind of watch mode or automatic image re-build.

### On development

```
go run ./cmd/api/
```

## Test

```
go test ./...
```

> Docker required: This project uses Testcontainers for integration tests (See ./integration_test directory)

### Examples

The directory `./examples` has http files with request functional examples that could be runned from IDE (Requires a HTTP client extension).
