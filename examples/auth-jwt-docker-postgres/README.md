# Auth + JWT + Docker + Postgres Example

This example extends the basic JWT authentication by integrating with a PostgreSQL database using GORM, all orchestrated with Docker Compose. It demonstrates a more realistic setup for a backend API with user registration, login, and protected routes.

## Features

-   User registration (`POST /auth/register`) with password hashing.
-   User login (`POST /auth/login`) returning a JWT.
-   Protected profile endpoint (`GET /profile`) requiring a valid JWT.
-   Admin-only protected endpoint (`GET /admin/dashboard`) demonstrating role-based authorization.
-   User data persisted in a PostgreSQL database.
-   Database schema migration using GORM `AutoMigrate`.
-   All services (Go API and PostgreSQL) managed via `docker-compose.yml`.
-   Environment-based configuration for database connection and JWT secret.

## Prerequisites

-   Docker and Docker Compose installed.
-   `jq` for parsing JSON responses in test commands.

## How to Run

1.  Navigate to the example directory:
    ```bash
    cd examples/auth-jwt-docker-postgres
    ```
2.  Build and start the services using Docker Compose:
    ```bash
    docker-compose up --build -d
    ```
    This will:
    -   Build the Go API service (`api`) from the `Dockerfile`.
    -   Start a PostgreSQL database service (`db`).
    -   The Go application will connect to the database, perform auto-migration, and start on port `8082`.

3.  Verify the services are running:
    ```bash
    docker-compose ps
    ```
    You should see `api` and `db` services in a healthy state.

## How to Test

The API service will be accessible at `http://localhost:8082`.

### 1. Register a New User (Happy Path)

```bash
curl -v -X POST "http://localhost:8082/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","email":"newuser@example.com","password":"securepass"}'
```
**Expected Response (201 Created):**
```json
{
    "message": "User registered successfully",
    "username": "newuser"
}
```

### 2. Login as the New User (Happy Path)

```bash
LOGIN_RESPONSE=$(curl -s -X POST "http://localhost:8082/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","password":"securepass"}')

echo "$LOGIN_RESPONSE"

ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token')
echo "Extracted Access Token: $ACCESS_TOKEN"
```
**Expected Response (200 OK):**
```json
{
    "access_token": "eyJ...",
    "token_type": "bearer",
    "expires_in": 900
}
```

### 3. Access Profile with New User's Token (Happy Path)

```bash
curl -v -X GET "http://localhost:8082/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```
**Expected Response (200 OK):**
```json
{
    "message": "Welcome to your profile!",
    "username": "newuser",
    "user_role": "user",
    "user_email": "newuser@example.com",
    "claims": {
        "sub": "newuser",
        "exp": ...,
        "iat": ...,
        "role": "user",
        "email": "newuser@example.com"
    }
}
```

### 4. Login as Default Admin User (Happy Path)

A default admin user (`username: admin`, `password: adminpass`) is seeded into the database if it doesn't exist.

```bash
ADMIN_LOGIN_RESPONSE=$(curl -s -X POST "http://localhost:8082/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"adminpass"}')

echo "$ADMIN_LOGIN_RESPONSE"

ADMIN_ACCESS_TOKEN=$(echo "$ADMIN_LOGIN_RESPONSE" | jq -r '.access_token')
echo "Extracted Admin Access Token: $ADMIN_ACCESS_TOKEN"
```

### 5. Access Admin Dashboard with Admin Token (Happy Path)

```bash
curl -v -X GET "http://localhost:8082/admin/dashboard" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN"
```
**Expected Response (200 OK):**
```json
{
    "message": "Welcome to the admin dashboard!",
    "admin": "admin"
}
```

### 6. Access Admin Dashboard with Regular User Token (Error Path: Forbidden)

```bash
curl -v -X GET "http://localhost:8082/admin/dashboard" \
  -H "Authorization: Bearer $ACCESS_TOKEN" # Using the 'newuser' token
```
**Expected Response (403 Forbidden):**
```json
{
    "error": {
        "code": "forbidden",
        "message": "Admin access required"
    }
}
```

### 7. Registration Error: Duplicate Username/Email

```bash
curl -v -X POST "http://localhost:8082/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","email":"another@example.com","password":"pass"}'
```
**Expected Response (409 Conflict):**
```json
{
    "error": {
        "code": "user_exists",
        "message": "Username or email already registered"
    }
```

## Cleanup

To stop and remove the Docker containers and associated volumes:
```bash
docker-compose down -v
```

## Code Highlights

-   **`main.go`**:
    -   Loads `DATABASE_URL` and `JWT_SECRET` from environment variables, crucial for Docker integration.
    -   Connects to PostgreSQL using GORM.
    -   `db.AutoMigrate(&User{})` automatically creates the `users` table.
    -   Includes `/auth/register` endpoint for creating new users with hashed passwords.
    -   `/auth/login` authenticates against the database and issues JWTs.
    -   `/profile` is a protected endpoint.
    -   `/admin/dashboard` demonstrates a simple role-based authorization check within a handler.
-   **`docker-compose.yml`**: Defines two services: `api` (your Go app) and `db` (PostgreSQL). It sets up networking and environment variables.
-   **`Dockerfile`**: Builds the Go application into a lightweight Alpine-based image.
-   **`User` struct**: Defines the GORM model for users, including `gorm:"uniqueIndex"` for constraints and `json:"-"` to prevent password hash from being serialized.

This example showcases a robust, production-ready authentication setup using `goBastion Lib` with external services.
