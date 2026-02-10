# Bulls Lab Backend

This is the backend service for Bulls Lab, built with [Go](https://golang.org/) and the [Gin Web Framework](https://gin-gonic.com/). It follows a hexagonal architecture (Ports and Adapters) for better maintainability and testability.

## 🚀 Getting Started

### Prerequisites

- **Go**: Version `1.24.5` or higher.
  - Check your Go version:
    ```bash
    go version
    ```

### Installation

1.  **Clone the repository** (if you haven't already).
2.  **Download dependencies**:
    ```bash
    go mod download
    ```

### Development

To start the server locally:

```bash
go run cmd/main.go
```

The server will be available at `http://localhost:8080/`.

## 🛠️ Commands

| Command | Description |
|---------|-------------|
| `go run cmd/main.go` | Starts the backend server. |
| `go build -o server cmd/main.go` | Builds the application into an executable named `server`. |
| `go test ./...` | Runs all tests in the project. |
| `go fmt ./...` | Formats the Go source code. |

## 📂 Project Architecture

The project follows a **Hexagonal Architecture**:

```text
bulls-lab-be/
├── cmd/
│   └── main.go          # Application entry point & dependency injection
├── internal/
│   ├── core/
│   │   ├── domain/      # Domain models (e.g., User)
│   │   ├── ports/       # Interface definitions for repositories/services
│   │   └── services/    # Business logic implementation
│   └── adapters/
│       ├── handler/     # HTTP/API handlers (Gin)
│       └── repository/  # Data persistence implementation (Memory/DB)
├── go.mod               # Module dependencies
└── go.sum               # Dependency checksums
```

## 🔌 API Endpoints

- `GET /health` - Health check endpoint
- `GET /users/:id` - Get user by ID
- `POST /users` - Create a new user

## 🧰 Technologies Used

- **Go 1.24.5**: Programming language
- **Gin**: HTTP web framework
- **Hexagonal Architecture**: Architectural pattern
