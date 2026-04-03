# Go Auth App

A robust, modular authentication API built with Go and Gin.  
Features secure user registration, email verification, JWT authentication (with refresh/rotation), logout, password reset, role-based access control (RBAC), profile & admin endpoints, and strong test coverage.  
Designed for rapid customization, maintainability, and solid testing.

## 🚀 Features (Grouped by Area)

### User Authentication
- **Register:** Secure new user registration with immediate email verification flow
- **Login:** Authenticate users (only after email verified)
- **JWT Auth:** Issue, validate, and rotate access/refresh tokens for each login
- **Logout:** Server-side token invalidation (refresh token blacklisting)

### Email Verification & Account Activation
- **Verification flow:** Sends verification email on register (`/verify-email?token=...`)
- **Resend Verification:** Endpoint to resend verification email

### Password Management
- **Forgot Password:** Request a password reset email if forgotten
- **Reset Password:** Reset password via secure emailed token
- **Secure password hashing:** All user passwords with bcrypt

### Token Management & Security
- **Token rotation:** Refresh token invalidation/rotation (prevents reuse after logout/refresh)
- **Access/Refresh tokens:** Short-lived JWT access, long-lived refresh
- **Token required for protected endpoints**

### User & Admin Operations
- **Profile:** Authenticated users can get their own profile details
- **User listing:** Admins can list/view all users (`/users`), RBAC enforced
- **Role-based guards:** Simple and extensible role checks (user/admin)

### Architecture & Developer UX
- **Clean, modular codebase:** Controllers, services, repos, DTOs, middleware
- **.env file support:** Easy environment config for DB, mail, JWT secrets
- **Mocks and tests:** Full coverage for core features

---

## 🏁 Quickstart

### Prerequisites

- Go 1.18 or newer ([get Go](https://golang.org/dl/))
- Docker (optional, for Postgres database)

### Installation

```sh
git clone https://github.com/your-username/go-auth-app.git
cd go-auth-app
go mod tidy
cp .env.example .env   # Edit .env with your DB, JWT, and email config
```

### Configuration

Key `.env` variables:
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE` — Postgres database connection settings
- `JWT_SECRET` — secret key for signing JWTs
- `ADMIN_EMAIL`, `ADMIN_PASSWORD` — initial admin user credentials
- `SENDGRID_API_KEY` — SendGrid API key for sending emails
- `SENDGRID_EMAIL` — email address used as sender for outgoing mail

### Running the Server

```sh
go run main.go
```

API is available at: [http://localhost:8080](http://localhost:8080)

---

## 📚 API Endpoints Grouped by Feature

### User Registration & Email Verification

- **POST `/register`** — Register a new user  
  **Request:**
  ```json
  { "name": "Alice", "email": "alice@email.com", "password": "supersecure" }
  ```
  - On success: User receives a verification email.

- **GET `/verify-email?token=...`** — Verify user's email  
  - Link sent in the verification email. Activates account.

- **POST `/resend`** — Resend verification email  
  **Request:**
  ```json
  { "email": "alice@email.com" }
  ```

---

### Authentication & Token Management

- **POST `/login`** — Log in (only after email verified)  
  **Request:**
  ```json
  { "email": "alice@email.com", "password": "supersecure" }
  ```
  **Response:**
  ```json
  { "access_token": "JWT...", "refresh_token": "..." }
  ```

- **POST `/refresh-token`** — Refresh tokens  
  **Request:**
  ```json
  { "refresh_token": "existing_refresh_token" }
  ```
  *Returns new pair, invalidates previous refresh token.*

- **POST `/logout`** — Invalidate issued refresh token  
  **Headers:** `Authorization: Bearer <access_token>`  
  **Request:**
  ```json
  { "refresh_token": "the_token_to_invalidate" }
  ```

---

### Password Reset

- **POST `/forgot-password`** — Send a password reset email  
  **Request:**
  ```json
  { "email": "alice@email.com" }
  ```

- **POST `/reset-password`** — Reset forgotten password  
  **Request:**
  ```json
  { "token": "<password_reset_token>", "new_password": "yourNewPassword" }
  ```

---

### User & Admin

- **GET `/profile`** — Get current logged-in user's profile  
  - **Requires:** valid access token

- **GET `/users`** — List all users *(admin only)*  
  - **Requires:** admin's access token (`Authorization: Bearer ...`)

---

## 🧪 Running Tests

```sh
go test ./tests/...
```
- Extensive coverage for registration, login, token lifecycles, role guards, verification logic, password reset, and more.
- Uses mocks for services and repositories, isolating business logic.

---

## 📂 Project Structure

```
.
├── controllers/     # HTTP handlers for API endpoints
├── dto/            # Request/response objects and validation
├── models/         # GORM models (User, Token, etc.)
├── repositories/   # Data interface & logic
├── services/       # Business logic (auth, user, mail, etc.)
├── middleware/     # JWT & RBAC middleware
├── config/         # Environment, DB, mail, JWT settings
├── tests/          # Unit/mocks for all features
└── main.go         # Entry point
```

---

## 🤝 Contributing

1. Fork this repo
2. Create a feature branch (`git checkout -b feat/your-feature`)
3. Commit and push your changes
4. Open a Pull Request!

## 📄 License

MIT License

---
