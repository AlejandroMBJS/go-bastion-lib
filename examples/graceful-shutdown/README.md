# Graceful Shutdown Example

This example demonstrates how to implement graceful shutdown for a `goBastion Lib` application. Graceful shutdown allows the server to finish processing ongoing requests within a specified timeout period before shutting down, preventing abrupt disconnections and data loss for active clients.

## Features

-   A `GET /long-task` endpoint that simulates a long-running operation (e.g., 5 seconds).
-   A `GET /health` endpoint to check server status.
-   The application listens for `SIGINT` (Ctrl+C) and `SIGTERM` signals.
-   Upon receiving a shutdown signal, the server attempts to shut down gracefully, allowing active requests to complete within a 10-second timeout.

## How to Run

1.  Navigate to the example directory:
    ```bash
    cd examples/graceful-shutdown
    ```
2.  Run the application:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8089`.

## How to Test

The API service will be accessible at `http://localhost:8089`.

### 1. Initiate a Long-Running Task

Open a **new terminal window** and execute a request to the `/long-task` endpoint. This request will take 5 seconds to complete.
```bash
curl http://localhost:8089/long-task
```
**Expected Response (after ~5 seconds, 200 OK):**
```json
{
    "message": "Long task completed"
}
```

### 2. Trigger Graceful Shutdown While Task is Running

While the `curl` command from step 1 is still executing (i.e., within the 5-second delay), go back to the terminal where your `go run main.go` command is running and press `Ctrl+C`.

**Expected Behavior:**

-   The server will log: `Shutting down server...`
-   The `curl` command in the other terminal will **still complete successfully** after its 5-second delay.
-   After the `curl` command finishes, the server will log: `Server exited gracefully.`
-   If the `curl` command took longer than the 10-second shutdown timeout, the server would log `Server forced to shutdown: context deadline exceeded`.

This demonstrates that the server waited for the active request to finish before fully terminating.

### 3. Test Health Endpoint

```bash
curl http://localhost:8089/health
```
**Expected Response (200 OK):**
```json
{
    "status": "ok"
}
```

## Code Highlights

-   **`main.go`**:
    -   The `app.Run()` method is called within a goroutine to allow the main goroutine to listen for OS signals.
    -   `signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)`: Configures the application to listen for interrupt (Ctrl+C) and termination signals.
    -   `<-quit`: Blocks the main goroutine until a signal is received.
    -   `context.WithTimeout(context.Background(), 10*time.Second)`: Creates a context with a 10-second deadline for the shutdown process.
    -   `app.Shutdown(ctx)`: This method (provided by `goBastion Lib`'s `App` abstraction) attempts to gracefully shut down the underlying `http.Server`. It waits for active connections to close or for the context deadline to be exceeded.

This example is vital for building resilient production applications that handle restarts and deployments without interrupting user experience.
