# CSRF Protection Example

This example demonstrates how to integrate and test CSRF (Cross-Site Request Forgery) protection using `goBastion Lib`'s built-in `middleware.CSRFMiddleware`. It uses the double-submit cookie pattern, where a token is sent in a cookie and must be echoed back in a custom HTTP header for "unsafe" requests.

## Features

-   `GET /form`: A public endpoint that, when accessed, causes the `CSRFMiddleware` to set a `csrf_token` cookie in the client's browser.
-   `POST /form/submit`: A protected endpoint that requires a valid `csrf_token` to be present in both the `csrf_token` cookie and the `X-CSRF-Token` header.
-   Proper `403 Forbidden` JSON response on CSRF validation failure.

## How to Run

1.  Navigate to the example directory:
    ```bash
    cd examples/csrf
    ```
2.  Run the application:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8083`.

## How to Test

### 1. Get the CSRF Token (Happy Path - GET Request)

First, make a GET request to `/form`. This will cause the server to set the `csrf_token` cookie. You need to extract this token for subsequent POST requests.

```bash
# Using curl, we need to capture the cookie.
# -c cookies.txt saves cookies to a file.
# -D headers.txt saves response headers to a file.
curl -v -c cookies.txt -D headers.txt http://localhost:8083/form
```
**Expected Response (200 OK):**
```json
{
    "message": "Load this form to get a CSRF cookie. Then extract X-CSRF-Token from the cookie.",
    "note": "Check your browser's cookies for 'csrf_token' after this request."
}
```
**Important:** After this command, inspect `cookies.txt` (or your browser's developer tools) to find the `csrf_token` value. It will look something like:
`#HttpOnly_localhost	FALSE	/	FALSE	0	csrf_token	a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2`
The long hex string is your CSRF token.

### 2. Submit Form with Valid CSRF Token (Happy Path - POST Request)

Now, use the extracted `csrf_token` in the `X-CSRF-Token` header and send the `cookies.txt` file back.

```bash
# Replace <YOUR_CSRF_TOKEN> with the actual token value from cookies.txt
# -b cookies.txt sends the cookies back to the server.
# -H "X-CSRF-Token: <YOUR_CSRF_TOKEN>" sends the token in the header.
CSRF_TOKEN="<YOUR_CSRF_TOKEN>" # e.g., a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2

curl -v -b cookies.txt -X POST "http://localhost:8083/form/submit" \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -d '{"message":"Hello from secure form!"}'
```
**Expected Response (200 OK):**
```json
{
    "status": "success",
    "message": "Form submitted successfully!",
    "data": "Hello from secure form!"
}
```

### 3. Submit Form with Missing CSRF Token (Error Path)

Attempt to submit the form without the `X-CSRF-Token` header.

```bash
curl -v -b cookies.txt -X POST "http://localhost:8083/form/submit" \
  -H "Content-Type: application/json" \
  -d '{"message":"Attempt without CSRF header"}'
```
**Expected Response (403 Forbidden):**
```json
{
    "error": "csrf_failed"
}
```

### 4. Submit Form with Invalid CSRF Token (Error Path)

Attempt to submit the form with an incorrect `X-CSRF-Token` header.

```bash
curl -v -b cookies.txt -X POST "http://localhost:8083/form/submit" \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: wrong-token" \
  -d '{"message":"Attempt with wrong CSRF token"}'
```
**Expected Response (403 Forbidden):**
```json
{
    "error": "csrf_failed"
}
```

## Code Highlights

-   **`main.go`**:
    -   `cfg.EnableCSRF = true` enables the CSRF middleware.
    -   `app.Use(middleware.CSRFMiddleware(middleware.DefaultCSRFConfig()))` applies the middleware globally.
    -   The `GET /form` endpoint serves as a way to trigger the middleware to set the `csrf_token` cookie.
    -   The `POST /form/submit` endpoint is protected by the CSRF middleware, which automatically validates the incoming token.

This example is crucial for web applications that rely on browser clients and cookie-based authentication.
