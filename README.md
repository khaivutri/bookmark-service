# 🔖 Bookmark Service

A lightweight Go-based REST API starter project designed to provide a simple and extensible foundation for bookmark-related applications.

The project currently exposes a health-check endpoint and follows a clean layered architecture, making it easy to extend with additional APIs and business logic.

---

## ✨ Overview

This project demonstrates a minimal production-style REST API using:

- **Go** as the server runtime
- **Gin** as the HTTP framework
- Environment-based configuration via environment variables or a `.env` file
- Layered architecture (Handler → Service → Model)
- Unit and integration testing

Although the current implementation is intentionally simple, the project structure is ready for future expansion.

---

## 🚀 Features

- Lightweight REST API server
- Health check endpoint
- Environment-based configuration
- Configurable service name and instance ID
- Automatic UUID generation for `INSTANCE_ID`
- Unit and integration tests
- Clean project structure

---

## 📁 Project Structure

```text
bookmark-service/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/                     # API engine and routing
│   ├── handler/                 # HTTP handlers
│   ├── model/                   # Shared data models
│   ├── service/                 # Business logic
│   └── integration_test/        # Integration tests
├── go.mod
└── go.sum
```

---

## 📋 Requirements

Before running the project, make sure you have:

- Go **1.26** or newer
- Git

---

# 🚀 Getting Started

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
```

> **Note**
>
> `INSTANCE_ID` is optional.
> If omitted or left empty, the application automatically generates a UUID when starting.

---

## 4. Use the Makefile

This project includes a `Makefile` for common development tasks.

### List available targets

```bash
make help
```

### Install dependencies

```bash
make deps
```

### Run the application

```bash
make run
```

Override environment values at runtime:

```bash
make run APP_PORT=9090 SERVICE_NAME=my_service
```

### Run swagger generation then start the app

```bash
make dev-run
```

This target first generates Swagger docs and then starts the application.

### Run with custom configuration

You can override the default environment variables directly from the command line.

```bash
make dev-run \
  APP_PORT=<your_custom_port> \
  SERVICE_NAME=<your_custom_service_name> \
  INSTANCE_ID=<your_uuid>
```

Example:

```bash
make dev-run \
  APP_PORT=9090 \
  SERVICE_NAME=bookmark-service-dev \
  INSTANCE_ID=550e8400-e29b-41d4-a716-446655440000
```

> **Note**
>
> - `APP_PORT` specifies the port on which the application listens.
> - `SERVICE_NAME` specifies the service name used by the application.
> - `INSTANCE_ID` **must be a valid UUID**. The application validates this value during startup and will fail to start if the provided UUID is invalid.

---

## 5. Verify the service

After the server starts successfully, open:

```text
http://localhost:8080
```

If you changed `APP_PORT`, replace `8080` with your configured port.

---

# 📡 API

## Health Check

| Method | Endpoint        | Description                        | Response |
| ------ | --------------- | ---------------------------------- | -------- |
| GET    | `/health-check` | Returns service health information | `200 OK` |

### Example Request

```http
GET /health-check HTTP/1.1
Host: localhost:8080
```

### Example Response

```json
{
  "message": "OK",
  "service_name": "bookmark_service",
  "instance_id": "c45f7d4f-f0d0-42dc-90d8-d5eb0f6dbe5e"
}
```

---

# ⚙️ Configuration

The application reads configuration from:

- Environment variables
- `.env` file (if available)

| Variable       | Default             | Description                  |
| -------------- | ------------------- | ---------------------------- |
| `APP_PORT`     | `8080`              | HTTP server port             |
| `SERVICE_NAME` | `bookmark_service`  | Service name                 |
| `INSTANCE_ID`  | Auto-generated UUID | Optional instance identifier |

---

# 🧪 Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test ./... -cover
```

---

# 📄 License

This project is intended as a starter template for building Go REST APIs and can be extended to fit your own requirements.
