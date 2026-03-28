# Go Auth App

A modern modular authentication API in Go and Gin.  
Supports secure registration, login, JWT authentication, refresh, logout, and simple RBAC (“admin”/“user”).  
Built for extensibility, testability, and fast project setup.

## 🚀 Features

- User registration and login endpoints
- Secure password hashing with bcrypt
- JWT-based authentication & token rotation (access/refresh tokens)
- Auto rotation & invalidation of used/expired refresh tokens
- Logout: deletes the refresh token (secure server-side, not just JWT expiry!)
- Middleware-protected routes (e.g., `/profile`, `/users`)
- Simple role-based guards (admin/user)
- Clean-layer architecture: repository, model, controller, service
- Extensively unit tested (controllers/services/repositories) in `/tests`
- Mocks for repo/business logic for reliable isolated testing
- `.env` support for easy configuration

## 🏁 Quickstart

### Prerequisites

- Go 1.18+ ([Download Go](https://golang.org/dl/))
- Docker (optional, for Postgres/local db)

### Installation

```sh
git clone https://github.com/your-username/go-auth-app.git
cd go-auth-app
go mod tidy
cp .env.example .env  # Edit .env for DB/JWT config as needed
```

### Configuration

- Edit `.env` as needed. Key variables:
  - `DATABASE_URL` (your Postgres data source name)
  - `JWT_SECRET`
  - `PORT` (default: 8080)

### Run the API

```sh
go run main.go
```
Server starts at: [http://localhost:8080](http://localhost:8080)

## 📚 API Endpoints

- `POST /register` — Register user  
  **Body:**
  ```json
  {
    "name": "Alice Smith",
    "email": "alice@email.com",
    "password": "supersecure"
  }
  ```

- `POST /login` — Obtain access + refresh tokens  
  **Body:**
  ```json
  {
    "email": "alice@email.com",
    "password": "supersecure"
  }
  ```
  **Success Response:**
  ```json
  {
    "access_token": "JWT...",
    "refresh_token": "..."
  }
  ```

- `POST /refresh-token` — Refresh tokens  
  **Body:**
  ```json
  {
    "refresh_token": "existing_refresh_token"
  }
  ```
  - On success: new `access_token` & `refresh_token` returned.
  - Old refresh token is invalidated (“rotation” for security).

- `POST /logout` — Logout & invalidate the provided refresh token  
  - Requires:  
    - `Authorization: Bearer <access_token>` header  
    - Body:
    ```json
    {
      "refresh_token": "the_token_to_invalidate"
    }
    ```
  - Logs user out of device/session.

- `GET /profile` — Fetch current user’s profile  
  - Requires: `Authorization: Bearer <access_token>`

- `GET /users` — (admin only)  
  - List all users.  
  - Requires: Authorization header from a user with `role: admin`.

## 🧪 Running Tests

Run all tests:
```sh
go test ./tests/...
```
- Unit and isolation tests are provided for mock repositories & controllers.
- Test covers: registration, login, failed/invalid login, refresh logic, protected endpoints, and admin guards.

## 📂 Project Structure

```
.
├── controllers/         # Handler/controller logic per endpoint
├── dto/                 # Request/response DTO structs
├── models/              # GORM models (User, RefreshToken, etc.)
├── repositories/        # Interface & storage logic (db repo, etc.)
├── services/            # Business logic (authentication, user, etc.)
├── middleware/          # JWT & role middleware
├── config/              # DB/JWT/env boot & settings
├── tests/               # Unit tests & mocking
└── main.go              # Entrypoint / main app bootstrap
```

## 🤝 Contributing

1. Fork this repository
2. Create a branch (`git checkout -b feat/my-feature`)
3. Commit & push your changes
4. Open a Pull Request!

## 📄 License

MIT License

---
