# Chirpy

Chirpy is a REST API that serves as the backend engine for a microblogging social network platform similar to X (formerly Twitter). It handles user authentication, post (chirp) management, and internal administration tasks.

## Features

- **User Account Management**: User registration, secure authentication via Argon2id password hashing, and credentials update.
- **Session Control**: Dual-token system using lightweight JSON Web Tokens (JWT) for access control and cryptographically secure random tokens for session refreshing and revocation.
- **Chirp Management**: Create, read, and delete short text posts (chirps) capped at 140 characters. Includes a lightweight automated profanity filter template for content moderation.
- **Query & Filters**: Advanced retrieval options for chirps, including filtering by author and sorting results (ascending/descending) in memory.
- **Webhooks & Subscriptions**: Integration endpoints to support premium feature upgrades (e.g., Chirpy Red membership) via external secure webhook signals protected by API key validation.
- **Admin Utilities**: Operational metrics monitoring (tracking page hits via atomic counters) and developer-mode state resetting capabilities.

## Tech Stack

- **Language**: Go (Golang)
- **Database**: PostgreSQL
- **SQL Generator**: SQLC (for type-safe SQL-to-Go generation)
- **Database Migrations**: Goose

## Prerequisites

Ensure you have the following installed on your local machine:

- Go (1.21 or higher)
- PostgreSQL instance running

## Getting Started

### 1. Environment Configuration

Create a `.env` file in the root directory of the project and populate the following variables:

```env
DB_URL="postgres://username:password@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="your_secure_jwt_signing_secret"
POLKA_KEY="your_external_webhook_api_key"

```

### 2. Run Migrations

Apply the database schema up to the latest version using Goose:

```bash
goose -dir db/migrations postgres "postgres://username:password@localhost:5432/chirpy?sslmode=disable" up

```

### 3. Build and Run the Server

Execute the entry point binary from the repository root:

```bash
go run cmd/chirpy/main.go

```

The application will start serving the API endpoints and local static files on port `8080`.

## API Endpoints

### Public / Utility

* `GET /api/healthz` - Service readiness check.

### Authentication & Users

* `POST /api/users` - Register a new account.
* `POST /api/login` - Authenticate and retrieve access/refresh tokens.
* `PUT /api/users` - Update authenticated user email and password (JWT required).
* `POST /api/refresh` - Request a new access token using a valid refresh token.
* `POST /api/revoke` - Log out by invalidating a refresh token.

### Chirps

* `POST /api/chirps` - Publish a new chirp (JWT required).
* `GET /api/chirps` - List all chirps (Supports optional `?author_id=` and `?sort=desc` filters).
* `GET /api/chirps/{chirpID}` - Retrieve a single chirp by its unique ID.
* `DELETE /api/chirps/{chirpID}` - Delete an existing chirp (Ownership verified via JWT).

### Internal / Platform

* `POST /api/polka/webhooks` - Secure external webhook for user subscription upgrades (API Key required).
* `GET /admin/metrics` - HTML monitoring view displaying total hits.
* `POST /admin/reset` - Truncate user records and reset metrics counters (Restricted to `dev` platform).