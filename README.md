# Bookmark Service

A lightweight Go REST API for creating short-lived bookmark links backed by Redis.

The service exposes health-check, short-link creation, and redirect endpoints. It follows a clean layered architecture so HTTP handling, business logic, repository access, and shared packages stay separated and easy to test.

---

## Overview

This project demonstrates a production-style REST API using:

- Go as the server runtime
- Gin as the HTTP framework
- Redis as the URL storage, and both Redis and PostgreSQL as dependencies checked by health checks
- Environment-based configuration via environment variables or a `.env` file
- Layered architecture: Handler -> Service -> Repository -> Package
- Structured error logging with zerolog
- Unit, package, repository, handler, and integration tests
- Docker and Docker Compose for local development

---

## Features

- Lightweight REST API server
- Health check endpoint with Redis and PostgreSQL dependency status
- User registration endpoint with request validation
- URL shortening endpoint with TTL support
- Redirect endpoint for generated short codes
- Random alphanumeric short-code generation
- Environment-based configuration
- Configurable app port, service name, instance ID, log level, and Redis connection
- Swagger documentation
- Dockerfile and Docker Compose setup
- Makefile shortcuts for common development tasks

---

## Project Structure

```text
bookmark-service/
├── cmd/
│   └── api/
│       └── main.go                  # Application entry point
├── internal/
│   ├── api/                         # API engine and route wiring
│   ├── handler/                     # HTTP handlers
│   │   └── v1/                      # Versioned link and user handlers with DTOs
│   ├── model/                       # Shared response models
│   ├── repository/                  # Database and Redis persistence adapters
│   ├── service/                     # Business logic
│   ├── integration_test/            # Endpoint integration tests
│   └── test/                        # Test fixtures and helpers
├── pkg/
|   ├── dbutils/                     # Handle datbase error
│   ├── logger/                      # Logging configuration
│   ├── redis/                       # Redis client, config, and test helper
|   ├── response                     # handle common error
│   ├── sqldb/                       # Postgres database client and config
│   └── utils/                       # Shared utilities
|   └── validation/                  # Custom validation
├── Dockerfile
├── docker-compose.yaml
├── Makefile
├── go.mod
└── go.sum
```

---

## Requirements

Before running the project, make sure you have:

- Go 1.26 or newer
- Redis
- Postgres
- Docker and Docker Compose, optional but recommended
- Git
- `make`, optional but recommended

---

# Getting Started

## 1. Clone the repository

```bash
git clone https://github.com/khaivutri/bookmark-service.git
cd bookmark-service
```

---

## 2. Install dependencies

Download all Go modules.

```bash
go mod download
```

Or use the Makefile:

```bash
make deps
```

---

## 3. Configure environment variables

You can configure the application in one of two ways:

- Create a `.env` file in the project root.
- Export environment variables directly from your terminal.

Example:

```env
APP_PORT=8080
SERVICE_NAME=bookmark_service
INSTANCE_ID=
LOG_LEVEL=info
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
DB_HOST=localhost
DB_USER=admin
DB_PASS=admin
DB_NAME=bookmark
DB_PORT=5432
```

`INSTANCE_ID` is optional. If omitted or left empty, the application automatically generates a UUID when starting.

When running the app inside Docker Compose, use the service hostnames:

```env
REDIS_ADDRESS=redis:6379
DB_HOST=postgres
DB_USER=admin
DB_PASS=admin
DB_NAME=bookmark
DB_PORT=5432
```

---

## 4. Start Redis and Postgres

Start Redis and Postgres with Docker Compose:

```bash
make docker-up
```

If you only need Redis for a local health-check run, use:

```bash
make docker-redis
```

Or run Redis and Postgres directly with Docker:

```bash
docker run --rm --name redis -p 6379:6379 redis:alpine

docker run --rm --name postgres -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=bookmark -p 5432:5432 postgres:alpine
```

If a container named `redis` already exists, remove it before starting a new one:

```bash
docker rm -f redis
```

---

## 5. Run the application locally

```bash
make run
```

Override environment values at runtime:

```bash
make run APP_PORT=9090 SERVICE_NAME=bookmark-service-dev LOG_LEVEL=debug
```

Run Swagger generation before starting the app:

```bash
make dev-run
```

---

## 6. Run with Docker Compose

Build and start both Redis and the bookmark service:

```bash
make docker-up
```

Follow logs:

```bash
make docker-logs
```

Stop services:

```bash
make docker-down
```

Equivalent Docker Compose commands:

```bash
docker-compose up --build
docker-compose logs -f
docker-compose down
```

---

## 7. Verify the service

After the server starts successfully, verify the health endpoint:

```bash
curl http://localhost:8080/health-check
```

If you changed `APP_PORT`, replace `8080` with your configured port.

---

# API

## Health Check

| Method | Endpoint        | Description                                  | Success Response |
| ------ | --------------- | -------------------------------------------- | ---------------- |
| GET    | `/health-check` | Returns service health and dependency status | `200 OK`         |

When Redis or PostgreSQL is unavailable, the endpoint returns `503 Service Unavailable` with the corresponding dependency status set to `DOWN` and `message` set to `DEGRADED`.

### Example Request

```http
GET /health-check HTTP/1.1
Host: localhost:8080
```

### Example Healthy Response

```json
{
  "message": "OK",
  "service_name": "bookmark_service",
  "instance_id": "c45f7d4f-f0d0-42dc-90d8-d5eb0f6dbe5e",
  "dependency": {
    "redis": "UP",
    "postgres": "UP"
  }
}
```

### Example Degraded Response (e.g., PostgreSQL Down)

```json
{
  "message": "DEGRADED",
  "service_name": "bookmark_service",
  "instance_id": "c45f7d4f-f0d0-42dc-90d8-d5eb0f6dbe5e",
  "dependency": {
    "redis": "UP",
    "postgres": "DOWN"
  }
}
```

---

## User Registration

| Method | Endpoint             | Description                                    | Success Response |
| ------ | -------------------- | ---------------------------------------------- | ---------------- |
| POST   | `/v1/users/register` | Registers a new user and returns user metadata | `201 Created`    |

The register endpoint validates the incoming JSON request and returns detailed error responses for invalid input or duplicate fields.

### Example Request

```http
POST /v1/users/register HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "username": "johndoe",
  "display_name": "John Doe",
  "password": "Password123@",
  "email": "john.doe@example.com"
}
```

### Example Success Response

```json
{
  "data": {
    "id": "7d80f755-7dce-4c95-b8bf-75bb8e240ef2",
    "username": "johndoe",
    "display_name": "John Doe",
    "email": "john.doe@example.com",
    "created_at": "2026-07-22T11:00:00Z",
    "updated_at": "2026-07-22T11:00:00Z"
  },
  "message": "User registered successfully!"
}
```

### Example Validation Error Response

```json
{
  "message": "Invalid input",
  "details": ["username is invalid (min)", "password is invalid (password)"]
}
```

### Example Duplicate Field Responses

```json
{
  "message": "username already exists"
}
```

```json
{
  "message": "email already exists"
}
```

### Example Internal Error Response

```json
{
  "message": "Processing error"
}
```

---

## Create Short Link

| Method | Endpoint            | Description                          | Success Response |
| ------ | ------------------- | ------------------------------------ | ---------------- |
| POST   | `/v1/links/shorten` | Creates a short code for a given URL | `200 OK`         |

### Request Body

| Field | Type   | Required | Validation           | Description                |
| ----- | ------ | -------- | -------------------- | -------------------------- |
| `url` | string | Yes      | Must be a valid URL  | Original URL to shorten    |
| `exp` | int64  | Yes      | Must be at least `5` | Expiration time in seconds |

### Example Request

```bash
curl -X POST http://localhost:8080/v1/links/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","exp":60}'
```

### Example Response

```json
{
  "code": "AbCDeFK",
  "message": "Shorten URL generated successfully!"
}
```

---

## Redirect

| Method | Endpoint                    | Description                              | Success Response |
| ------ | --------------------------- | ---------------------------------------- | ---------------- |
| GET    | `/v1/links/redirect/{code}` | Redirects a short code to the stored URL | `302 Found`      |

### Example Request

```bash
curl -i http://localhost:8080/v1/links/redirect/AbCDeFK
```

If the code exists, the response includes a `Location` header pointing to the original URL.

If the code does not exist or has expired, the endpoint returns:

```json
{
  "error": "Code not found"
}
```

---

# Swagger

Generate Swagger docs:

```bash
make swagger
```

Run the app and open:

```text
http://localhost:8080/swagger/index.html
```

---

# Configuration

The application reads configuration from:

- Environment variables
- `.env` file, if available

| Variable         | Default             | Description                                   |
| ---------------- | ------------------- | --------------------------------------------- |
| `APP_PORT`       | `8080`              | HTTP server port                              |
| `SERVICE_NAME`   | `bookmark_service`  | Service name returned by the health-check API |
| `INSTANCE_ID`    | Auto-generated UUID | Optional instance identifier                  |
| `LOG_LEVEL`      | `info`              | Log level passed to zerolog                   |
| `REDIS_ADDRESS`  | `localhost:6379`    | Redis host and port                           |
| `REDIS_PASSWORD` | empty               | Redis password                                |
| `REDIS_DB`       | `0`                 | Redis database number                         |

`INSTANCE_ID` must be a valid UUID when explicitly provided. Invalid values cause startup to fail.

---

# Makefile Commands

## Makefile Commands

| Command               | Description                                          |
| --------------------- | ---------------------------------------------------- |
| `make help`           | Show available Makefile targets                      |
| `make deps`           | Download Go module dependencies                      |
| `make tidy`           | Clean up Go module dependencies                      |
| `make test`           | Run all tests with coverage report                   |
| `make run`            | Run the application locally                          |
| `make clean`          | Remove build artifacts and coverage files            |
| `make swagger`        | Generate Swagger documentation                       |
| `make dev-run`        | Generate Swagger docs and run the application        |
| `make docker-up`      | Build and start services with Docker Compose         |
| `make docker-down`    | Stop Docker Compose services                         |
| `make docker-logs`    | Follow Docker Compose logs                           |
| `make docker-redis`   | Start only Redis with Docker Compose                 |
| `make docker-test`    | Run tests inside Docker and generate coverage report |
| `make docker-build`   | Build the Docker image                               |
| `make docker-login`   | Log in to a Docker registry                          |
| `make docker-release` | Push the Docker image to the registry                |

---

# Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage through the Makefile:

```bash
make test
```

The test suite includes:

- Handler tests
- Service tests
- Repository tests
- Redis client/config/mock tests
- Logger level tests
- Code generator tests
- Integration tests for health-check and short-link endpoints

---

# CI Pipeline

This project uses GitHub Actions to automate testing and deployment.

**Triggers:**

- Pull requests targeting the `main` branch
- Pushes to the `main` branch or tags matching `v*.*.*`

**Pipeline steps:**

1. **Checkout code** – fetch the source code from the repository
2. **Set up Go** – install Go using the version declared in `go.mod`, with caching enabled for faster builds
3. **Run unit tests** – run `make test` to execute the test suite
4. **Log in to Docker Hub** _(push only)_ – authenticate with Docker Hub using stored secrets
5. **Build & push Docker image** _(push only)_ – build the Docker image and push it to Docker Hub with the tag `khaivutri/shorten_link:latest`

**Code security:**
The repository is integrated with **SonarQube** for static code analysis, helping detect security vulnerabilities, code smells, and coding convention violations early, before code is merged.

> Note: Unit tests always run, while Docker build & push only run on pushes to `main`. Pull requests only run checkout, Go setup, and unit tests.

---

# License

This project is intended as a starter service for building Go REST APIs and can be extended to fit your own requirements.
