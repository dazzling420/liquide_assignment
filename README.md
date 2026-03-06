# Liquide Assignment

A backend application built with Go that provides an API for user authentication and placing buy orders, utilizing MongoDB for data persistence and Redis for caching or fast storage.

## Tech Stack

- **Language:** Go (1.25.3)
- **Web Framework:** chi (v1.5.5)
- **Database:**
  - MongoDB (via `mongo-driver`)
  - Redis (via `rueidis`)
- **Logging:** Uber Zap (v1.27.1)
- **Authentication:** JWT (`golang-jwt/jwt/v5`)
- **Utilities:** `google/uuid`, `golang.org/x/crypto`, `gopkg.in/yaml.v3`

## Project Structure

- `cmd/server/`: Contains the entry point of the application (`main.go`).
- `config/`: Contains YAML configuration files (e.g., `testConfig.yml`).
- `internal/`: Contains core application code:
  - `config/`: Configuration models and loading logic.
  - `logger/`: Logging setup.
  - `models/`: Data models defined for the application.
  - `repository/`: Data access layers.
  - `server/`: Server startup and shutdown routines.
  - `service/`: Core business logic.
  - `storage/`: Database connection and interface implementations (MongoDB, Redis).

## Setup and Configuration

The application loads configuration from the `config/` directory based on the `-env` flag. By Default the env is set to "test"

### Running the Application

To run the application, jump to the `cmd/server` directory and execute:

```bash
cd cmd/server
go run main.go -env test
```

The server listens on port `8000` by default.

## API Endpoints

The following REST API endpoints are exposed. 

### Authentication

- **POST `/signup`**: Register a new user.
  - Body: `{"username": "your_username", "password": "your_password"}`
- **POST `/login`**: Authenticate a user and receive a JWT token.
  - Body: `{"username": "your_username", "password": "your_password"}`
  - Returns: A JWT token to be used in the `Authorization` header.

### Orders

Requires the `Authorization: Bearer <token>` header for access.

- **POST `/buy`**: Place a new buy order.
  - Body: `{"price": 100, "quantity": 10}`
- **GET `/orderbook`**: Retrieve the current order book and view existing orders.
