# Hello World + Custom 404 Handler Example

This example demonstrates the most basic usage of `goBastion Lib` by setting up a "Hello, World!" endpoint and a custom 404 Not Found handler.

## Features

- Minimal API setup.
- `GET /` endpoint returning "Hello, World!".
- Custom JSON response for 404 Not Found errors.

## How to Run

1. Navigate to the example directory:
   ```bash
   cd examples/hello-404
   ```
2. Run the application:
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:8080`.

## How to Test

### 1. Hello World Endpoint (Happy Path)

Open your browser or use `curl`:
```bash
curl http://localhost:8080/
```
**Expected Response (200 OK):**
```json
{
    "message": "Hello, World!"
}
```

### 2. Custom 404 Handler (Error Path)

Try to access a non-existent endpoint:
```bash
curl http://localhost:8080/non-existent-path
```
**Expected Response (404 Not Found):**
```json
{
    "error": "not_found",
    "message": "The requested resource could not be found on this server."
}
```

## Code Highlights

```go
// --- Custom 404 Not Found Handler ---
// This demonstrates how to override the default 404 JSON response.
r.SetNotFoundHandler(func(ctx *router.Context) {
	ctx.JSON(http.StatusNotFound, map[string]string{
		"error":   "not_found",
		"message": "The requested resource could not be found on this server.",
	})
})

// --- Hello World Endpoint ---
r.GET("/", func(ctx *router.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"message": "Hello, World!",
	})
})
```
This example showcases how easy it is to define routes and customize framework behavior like the 404 handler.
