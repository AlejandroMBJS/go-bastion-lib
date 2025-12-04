# Unit Testing & Integration Testing Example

This example demonstrates how to effectively test `goBastion Lib` applications at both the unit and integration levels using Go's standard `testing` package and `net/http/httptest`, along with `stretchr/testify` for assertions.

## Features

-   **Unit Tests for Handlers**: Shows how to test individual handler functions in isolation by manually creating a `router.Context` and `httptest.ResponseRecorder`.
-   **Integration Tests for Router**: Demonstrates how to test the entire routing and middleware chain by simulating HTTP requests against the `App`'s `Handler()`.
-   Covers testing for successful responses, error conditions, and path parameter extraction.
-   Uses `github.com/stretchr/testify/assert` for clear and concise assertions.

## How to Run

1.  Navigate to the example directory:
    ```bash
    cd examples/testing
    ```
2.  Run the tests:
    ```bash
    go test -v ./...
    ```
    The `-v` flag provides verbose output, showing each test case.

## How to Test (Expected Output)

When you run `go test -v ./...`, you should see output similar to this, indicating all tests passed:

```
=== RUN   TestPingHandler
--- PASS: TestPingHandler (0.00s)
=== RUN   TestListUsersHandler
--- PASS: TestListUsersHandler (0.00s)
=== RUN   TestCreateUserHandler_Success
--- PASS: TestCreateUserHandler_Success (0.00s)
=== RUN   TestCreateUserHandler_InvalidInput
--- PASS: TestCreateUserHandler_InvalidInput (0.00s)
=== RUN   TestGetUserHandler_Success
--- PASS: TestGetUserHandler_Success (0.00s)
=== RUN   TestGetUserHandler_NotFound
--- PASS: TestGetUserHandler_NotFound (0.00s)
=== RUN   TestRouterIntegration
=== RUN   TestRouterIntegration/GET_/ping
--- PASS: TestRouterIntegration/GET_/ping (0.00s)
=== RUN   TestRouterIntegration/GET_/users
--- PASS: TestRouterIntegration/GET_/users (0.00s)
=== RUN   TestRouterIntegration/POST_/users
--- PASS: TestRouterIntegration/POST_/users (0.00s)
=== RUN   TestRouterIntegration/GET_/users/:id
--- PASS: TestRouterIntegration/GET_/users/:id (0.00s)
=== RUN   TestRouterIntegration/GET_/nonexistent
--- PASS: TestRouterIntegration/GET_/nonexistent (0.00s)
--- PASS: TestRouterIntegration (0.00s)
PASS
ok      github.com/alejandrombjs/go-bastion-lib/examples/testing    0.008s
```

## Code Highlights

-   **`main.go`**: Contains a simple `goBastion Lib` application with handlers for `/ping`, `/users` (list/create), and `/users/:id` (get). This is the code under test.
-   **`main_test.go`**:
    -   **Unit Tests (`TestPingHandler`, `TestListUsersHandler`, etc.)**:
        -   Each test function creates an `httptest.ResponseRecorder` and `httptest.NewRequest` to simulate an HTTP request.
        -   A `router.Context` is manually created and passed to the handler function.
        -   Assertions are made on `w.Code` (HTTP status code) and `w.Body` (response body).
        -   A custom `SetParam` method is added to `router.Context` (for testing purposes) to simulate path parameters.
    -   **Integration Test (`TestRouterIntegration`)**:
        -   Sets up a full `bastion.App` instance, including global middlewares and all routes.
        -   Uses `httptest.NewRequest` and `httptest.NewRecorder` to send requests to `r.Handler().ServeHTTP(rec, req)`.
        -   This tests the entire request lifecycle, including middleware execution and route matching.
    -   **`github.com/stretchr/testify/assert`**: Provides a rich set of assertion functions that make tests more readable and robust.

This example provides a clear blueprint for ensuring the quality and correctness of your `goBastion Lib` APIs through automated testing.
