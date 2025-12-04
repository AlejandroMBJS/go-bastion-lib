# Docker + Nginx Reverse Proxy Example

This example demonstrates how to deploy a `goBastion Lib` application behind an Nginx reverse proxy using Docker Compose. This is a common setup for production environments, providing benefits like SSL/TLS termination, load balancing, and static file serving (though not fully demonstrated here).

## Features

-   A simple `goBastion Lib` API (`/api/hello`, `/api/headers`).
-   Nginx configured as a reverse proxy to forward requests to the Go API.
-   Docker Compose to orchestrate both the Nginx and Go API services.
-   Demonstrates how Nginx adds/modifies headers (e.g., `X-Real-IP`, `X-Forwarded-For`).

## Prerequisites

-   Docker and Docker Compose installed.

## How to Run

1.  Navigate to the example directory:
    ```bash
    cd examples/docker-nginx
    ```
2.  Build and start the services using Docker Compose:
    ```bash
    docker-compose up --build -d
    ```
    This will:
    -   Build the Go API service (`api`) from its `Dockerfile`.
    -   Start an Nginx service (`nginx`) and configure it to proxy requests to the `api` service.
    -   Nginx will be accessible on `http://localhost:80`.

3.  Verify the services are running:
    ```bash
    docker-compose ps
    ```
    You should see `nginx` and `api` services in a healthy state.

## How to Test

The API will be accessible through Nginx at `http://localhost`.

### 1. Access the Hello Endpoint via Nginx

```bash
curl http://localhost/api/hello
```
**Expected Response (200 OK):**
```json
{
    "message": "Hello from Go Bastion App!",
    "served_by": "Port 8080"
}
```
This shows that Nginx successfully proxied the request to the Go application.

### 2. Inspect Headers Received by the Go App

```bash
curl http://localhost/api/headers
```
**Expected Response (200 OK):**
You will see a JSON object containing all headers received by the Go application. Notice the `X-Real-Ip` and `X-Forwarded-For` headers added by Nginx, which are crucial for the Go app to correctly identify the client's original IP address (e.g., for rate limiting).
```json
{
    "message": "Request headers received by Go Bastion App",
    "headers": {
        "Accept": "*/*",
        "Host": "localhost",
        "User-Agent": "curl/...",
        "X-Forwarded-For": "172.18.0.1", # Example IP from Docker's internal network
        "X-Forwarded-Host": "localhost",
        "X-Forwarded-Port": "80",
        "X-Forwarded-Proto": "http",
        "X-Real-Ip": "172.18.0.1" # Example IP from Docker's internal network
    }
}
```

## Cleanup

To stop and remove the Docker containers:
```bash
docker-compose down
```

## Code Highlights

-   **`main.go`**: A basic `goBastion Lib` app that listens on a configurable port (default 8080) and exposes two endpoints.
-   **`Dockerfile`**: Builds the Go application into a Docker image.
-   **`nginx.conf`**: Configures Nginx to:
    -   Listen on port 80.
    -   `proxy_pass http://api:8080;`: Forwards all incoming requests to the `api` service (our Go app) running on port 8080 within the Docker network.
    -   `proxy_set_header ...`: Important headers are set to pass client information (like original IP) from Nginx to the backend Go application.
-   **`docker-compose.yml`**: Defines two services: `nginx` and `api`. It sets up the necessary port mappings, volume mounts for Nginx config, and networking.

This example is fundamental for understanding how to deploy `goBastion Lib` applications in a production-like environment.
