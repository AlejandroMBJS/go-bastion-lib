# JWT Authentication Example

This example demonstrates how to implement JWT (JSON Web Token) based authentication using `goBastion Lib`. It includes a login endpoint to obtain a JWT and a protected endpoint that requires a valid JWT for access.

## Features

-   `POST /auth/login`: Authenticates a user with hardcoded credentials and returns an `access_token`.
-   `GET /profile`: A protected endpoint that returns user information only if a valid JWT is provided in the `Authorization` header.
-   Integration of `middleware.JWTAuth` for protecting routes.
-   Reading JWT claims from the `router.Context`.

## How to Run

1.  Navigate to the example directory:
    ```bash
    cd examples/auth-jwt
    ```
2.  Run the application:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8081`.

## How to Test

### 1. Login to Get JWT Token (Happy Path)

First, obtain an `access_token` by logging in. The example uses hardcoded credentials: `username: testuser`, `password: testpass`.

```bash
LOGIN_RESPONSE=$(curl -s -X POST "http://localhost:8081/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}')

echo "$LOGIN_RESPONSE"

# Extract the access token for subsequent requests
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

### 2. Access Protected Endpoint (Happy Path)

Use the `ACCESS_TOKEN` obtained from the login step to access the protected `/profile` endpoint.

```bash
curl -v -X GET "http://localhost:8081/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```
**Expected Response (200 OK):**
```json
{
    "message": "Welcome to your profile!",
    "username": "testuser",
    "user_role": "user",
    "claims": {
        "sub": "testuser",
        "exp": 1678886400,
        "iat": 1678885500,
        "role": "user"
    }
}
```

### 3. Access Protected Endpoint (Error Path: Missing JWT)

Attempt to access `/profile` without providing a JWT.

```bash
curl -v -X GET "http://localhost:8081/profile"
```
**Expected Response (401 Unauthorized):**
```json
{
    "error": "unauthorized"
}
```

### 4. Access Protected Endpoint (Error Path: Invalid JWT)

Attempt to access `/profile` with an invalid or malformed JWT.

```bash
curl -v -X GET "http://localhost:8081/profile" \
  -H "Authorization: Bearer invalid.jwt.token"
```
**Expected Response (401 Unauthorized):**
```json
{
    "error": "unauthorized"
}
```

### 5. Login (Error Path: Invalid Credentials)

Attempt to login with incorrect credentials.

```bash
curl -v -X POST "http://localhost:8081/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"wronguser","password":"wrongpass"}'
```
**Expected Response (401 Unauthorized):**
```json
{
    "error": "unauthorized",
    "message": "Invalid credentials"
}
```

## Code Highlights

The `main.go` demonstrates:

-   Configuring `bastion.Config` to enable JWT and set the `JWTSecret`.
-   Defining a `POST /auth/login` handler that uses `security.GenerateAccessToken` to create a JWT upon successful authentication.
-   Creating a protected route group using `r.Group("/")` and applying `middleware.JWTAuth(cfg.JWTSecret)` to it.
-   Accessing the decoded JWT claims within a protected handler via `ctx.Get("userClaims")`.

This example is a core building block for any API requiring user authentication.
