# Basic Example

This example demonstrates how to use the go-bastion-lib framework.

## Prerequisites

- Go 1.21 or later
- git

## Running the Example

1. Clone the repository:
   ```bash
   git clone https://github.com/alejandrombjs/go-bastion-lib.git
   cd go-bastion-lib/examples/basic
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

3. The server will start on `http://localhost:9876`

## Testing with cURL

### 1. Check Health Endpoint
```bash
curl http://localhost:9876/api/health
```

### 2. Login to Get JWT Token
```bash
curl -X POST http://localhost:9876/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'
```

### 3. Access Protected Endpoint (with JWT)
```bash
# Replace <TOKEN> with the access_token from the login response
curl http://localhost:9876/api/users \
  -H "Authorization: Bearer <TOKEN>"
```

### 4. Create New User
```bash
curl -X POST http://localhost:9876/api/users \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","email":"new@example.com","password":"securepassword123"}'
```

## Features Demonstrated

- ✅ JWT Authentication
- ✅ Rate Limiting (100 requests per minute)
- ✅ Security Headers
- ✅ Request ID tracking
- ✅ Structured Logging
- ✅ Panic Recovery
- ✅ Graceful Shutdown
- ✅ JSON Request/Response handling
- ✅ Router with groups
- ✅ Middleware chain

## Project Structure

The example shows:
- Global middleware registration
- Route grouping (`/api`, `/api/auth`, `/api/users`)
- Protected routes with JWT middleware
- Error handling helpers
- JSON request binding
- Password hashing
