# WebSocket Echo Server Example

This example demonstrates how to integrate WebSocket functionality into a `goBastion Lib` application. It sets up a simple echo WebSocket server that receives messages from a client and sends them back.

## Features

-   `GET /ws`: Endpoint that upgrades a standard HTTP connection to a WebSocket connection.
-   Echo functionality: Any message received over the WebSocket is immediately sent back to the client.
-   Uses `github.com/gorilla/websocket` for robust WebSocket handling.
-   Includes a regular HTTP `/health` endpoint for server status.

## Prerequisites

-   A WebSocket client for testing (e.g., a browser's developer console, `wscat`, or a simple HTML page).

## How to Run

1.  Navigate to the example directory:
    ```bash
    cd examples/websocket
    ```
2.  Run the application:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8086`.

## How to Test

The WebSocket endpoint will be at `ws://localhost:8086/ws`.

### 1. Test with a Browser's Developer Console

1.  Open your web browser (e.g., Chrome, Firefox).
2.  Open the Developer Tools (usually F12).
3.  Go to the "Console" tab.
4.  Execute the following JavaScript code:

    ```javascript
    const ws = new WebSocket("ws://localhost:8086/ws");

    ws.onopen = (event) => {
        console.log("WebSocket connection opened:", event);
        ws.send("Hello from browser!");
    };

    ws.onmessage = (event) => {
        console.log("Message from server:", event.data);
    };

    ws.onclose = (event) => {
        console.log("WebSocket connection closed:", event);
    };

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    // To send more messages:
    // ws.send("Another message!");

    // To close the connection:
    // ws.close();
    ```
    You should see "Hello from browser!" echoed back in the console.

### 2. Test with `wscat` (Node.js tool)

If you have Node.js installed, you can install `wscat`:
```bash
npm install -g wscat
```
Then connect and send messages:
```bash
wscat -c ws://localhost:8086/ws
```
Type a message and press Enter. The server should echo it back.
```
> Hello wscat!
< Hello wscat!
> How are you?
< How are you?
```

### 3. Test HTTP Health Endpoint

```bash
curl http://localhost:8086/health
```
**Expected Response (200 OK):**
```json
{
    "status": "ok",
    "type": "http"
}
```

## Code Highlights

-   **`main.go`**:
    -   `websocket.Upgrader`: Configured to handle the WebSocket handshake. `CheckOrigin` is set to `true` for demonstration purposes, but should be restricted in production.
    -   `upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request().Request, nil)`: This is the core call that upgrades the HTTP connection to a WebSocket.
    -   `conn.ReadMessage()` and `conn.WriteMessage()`: Used in a loop to continuously read from and write to the WebSocket connection, implementing the echo functionality.
    -   Error handling for connection upgrades and message read/write operations.

This example provides a foundation for building real-time features like chat applications, live dashboards, or notifications within your `goBastion Lib` API.
