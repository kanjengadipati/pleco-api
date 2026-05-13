# Troubleshooting

This page focuses on the most common `pleco-api` setup and runtime issues.

## The API will not start because configuration is invalid

Symptoms:

- The server exits immediately on boot
- Startup logs mention missing or invalid environment values

Checks:

- Confirm you copied the correct example env file
- Confirm `JWT_SECRET`, `DATABASE_URL`, `APP_BASE_URL`, and `FRONTEND_URL` are set for your environment
- Make sure boolean and numeric settings do not contain stray quotes or spaces

## Database connection failed

Symptoms:

- Startup logs show connection refused or authentication errors
- Migrations fail before the API boots

Checks:

- Verify `DATABASE_URL` points to the correct host, port, database, username, and password
- Confirm PostgreSQL is running and reachable from the API host
- Confirm the `sslmode` matches your environment

Local example:

```env
DATABASE_URL=postgresql://postgres:password@localhost:5432/auth_db?sslmode=disable
```

## CORS preflight fails from the dashboard or browser client

Symptoms:

- Browser console shows a CORS error
- Requests fail before hitting the handler

Checks:

- Confirm `CORS_ALLOWED_ORIGINS` includes the exact frontend origin
- Do not use a wildcard origin with credentialed auth requests
- Verify the scheme and port are correct

Example:

```env
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
```

## Login succeeds but refresh-based flows fail

Symptoms:

- Initial login works
- Refresh, logout-other-sessions, or protected browser flows fail afterward

Checks:

- Confirm the client is receiving the `pleco_refresh_token` HttpOnly cookie
- Confirm the frontend origin is allowed in `CORS_ALLOWED_ORIGINS`
- In production, make sure the API is served over HTTPS so secure cookies can work correctly

## Social login returns provider validation errors

Checks:

- Set `SOCIAL_ACTIVE_PROVIDERS` to the providers you actually want enabled
- Confirm the provider client ID matches the token audience
- For Facebook, also confirm `SOCIAL_FACEBOOK_CLIENT_SECRET` is configured when Facebook auth is enabled

## Email verification or password reset links point to the wrong host

Checks:

- Confirm `APP_BASE_URL` points to the API host
- Confirm `FRONTEND_URL` points to the browser app host
- Restart the process after changing environment variables

Why this happens:

Pleco uses those values when building verification and recovery links.

## AI investigation is unavailable or returns provider errors

Checks:

- Confirm `AI_ENABLED=true`
- Confirm `AI_PROVIDER` is one of `mock`, `ollama`, `openai`, or `gemini`
- If using Ollama, confirm `AI_BASE_URL` is reachable and the model is installed
- If using OpenAI or Gemini, confirm `AI_API_KEY` is present and valid

Quick local test setup:

```env
AI_ENABLED=true
AI_PROVIDER=mock
AI_MODEL=mock-model
```

## Rate limiting behaves differently between local and production

Why this happens:

Without Redis, Pleco falls back to in-memory stores that work for local single-instance development.

Checks:

- Configure `REDIS_URL` or `REDIS_HOST` and `REDIS_PORT` for multi-instance deployments
- Confirm Redis is reachable from the API container or host

## Docker Compose starts, but the API is unhealthy

Checks:

- Run `make docker-logs`
- Confirm Postgres finished initializing
- Confirm the Docker env file exists and matches the compose setup
- Check whether a stale local `.env.docker` is overriding expected values

## Newman or Postman tests fail locally

Checks:

- Start the API before running `make postman-all`
- Confirm the Postman environment file points at the local API
- Run `npm ci` if Newman dependencies are missing

## Migrations are out of sync or stuck

Checks:

- Run `make migrate-status`
- Use `make migrate-force VERSION=<number>` only when you understand the current database state
- Re-run `make migrate-up` after correcting the migration version

## Port 8080 is already in use

Checks:

- Stop the conflicting process
- Or set `PORT` to another value in your env file before starting the API

## Still stuck?

When opening an issue or asking for help, include:

- Pleco version or Git tag
- Go version
- How you are running the app: local Go, Docker Compose, or release binary
- Relevant env values with secrets removed
- The exact error message or failing endpoint
