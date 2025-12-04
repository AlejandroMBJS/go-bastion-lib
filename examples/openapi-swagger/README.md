# OpenAPI + Swagger Docs Example

This example demonstrates how to serve an OpenAPI (Swagger) specification file and integrate with Swagger UI using `goBastion Lib`. This allows you to provide interactive API documentation directly from your application.

## Features

-   `GET /openapi.yaml`: Serves the raw OpenAPI 3.0 specification file.
-   `GET /docs`: Redirects to a public Swagger UI instance, configured to load the `/openapi.yaml` from this server.
-   Includes a simple `/api/hello` endpoint that is documented in `openapi.yaml`.

## Prerequisites

-   Internet connection to access the public Swagger UI (petstore.swagger.io).

## How to Run

1.  Navigate to the example directory:
    ```bash
    cd examples/openapi-swagger
    ```
2.  Run the application:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8087`.

## How to Test

The API service will be accessible at `http://localhost:8087`.

### 1. Access the Raw OpenAPI Specification

Open your browser or use `curl`:
```bash
curl http://localhost:8087/openapi.yaml
```
**Expected Response (200 OK):**
The content of the `openapi.yaml` file will be returned.

### 2. Access the Interactive Swagger UI

Open your web browser and navigate to:
```
http://localhost:8087/docs
```
This will redirect you to `https://petstore.swagger.io/?url=http://localhost:8087/openapi.yaml`. The Swagger UI should load and display the documentation for the `/api/hello` endpoint.

You can then use the Swagger UI to interact with the `/api/hello` endpoint.

### 3. Test the Documented API Endpoint

```bash
curl http://localhost:8087/api/hello
```
**Expected Response (200 OK):**
```json
{
    "message": "Hello, World!"
}
```

```bash
curl "http://localhost:8087/api/hello?name=Alice"
```
**Expected Response (200 OK):**
```json
{
    "message": "Hello, Alice!"
}
```

## Code Highlights

-   **`main.go`**:
    -   `r.GET("/openapi.yaml", ...)`: Uses `http.ServeFile` to serve the `openapi.yaml` file directly from the application's file system.
    -   `r.GET("/docs", ...)`: Performs an `http.Redirect` to a public Swagger UI instance, passing the URL of our local `openapi.yaml` as a query parameter. This is a convenient way to get interactive docs without embedding the Swagger UI static files.
    -   The `/api/hello` endpoint is a simple example of an API that would be documented in the OpenAPI spec.
-   **`openapi.yaml`**: Contains the OpenAPI 3.0 specification for the `/api/hello` endpoint, detailing its parameters and responses.

This example is crucial for providing clear, interactive, and machine-readable documentation for your `goBastion Lib` APIs.
