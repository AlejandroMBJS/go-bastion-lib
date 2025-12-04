# File Upload and Serve Example

This example demonstrates how to handle file uploads using multipart/form-data and how to serve static files from a specified directory with `goBastion Lib`.

## Features

-   `POST /upload`: Accepts a file via multipart form data, saves it to a local `uploads/` directory.
-   `GET /files/:name`: Serves a previously uploaded file by its name.
-   Basic filename sanitization to prevent path traversal attacks.
-   Error handling for various upload and serve scenarios.

## Prerequisites

-   A local `uploads/` directory (will be created automatically if it doesn't exist).

## How to Run

1.  Navigate to the example directory:
    ```bash
    cd examples/file-upload
    ```
2.  Run the application:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8085`.

## How to Test

The API service will be accessible at `http://localhost:8085`.

### 1. Upload a File (Happy Path - POST /upload)

First, create a dummy file to upload.
```bash
echo "This is a test file content." > test_upload.txt
```

Now, upload it:
```bash
curl -v -X POST "http://localhost:8085/upload" \
  -F "file=@test_upload.txt"
```
**Expected Response (200 OK):**
```json
{
    "message": "File uploaded successfully",
    "filename": "test_upload.txt",
    "size": "27 bytes",
    "location": "./uploads/test_upload.txt"
}
```
You should also see `test_upload.txt` created inside the `examples/file-upload/uploads/` directory.

### 2. Serve the Uploaded File (Happy Path - GET /files/:name)

```bash
curl -v "http://localhost:8085/files/test_upload.txt"
```
**Expected Response (200 OK):**
The response body should contain the content of `test_upload.txt`:
```
This is a test file content.
```

### 3. Upload a File (Error Path: Missing File)

```bash
curl -v -X POST "http://localhost:8085/upload"
```
**Expected Response (400 Bad Request):**
```json
{
    "error": {
        "code": "missing_file",
        "message": "Error retrieving file from form: http: no such file"
    }
}
```

### 4. Serve a Non-Existent File (Error Path: Not Found)

```bash
curl -v "http://localhost:8085/files/non_existent_file.txt"
```
**Expected Response (404 Not Found):**
```json
{
    "error": {
        "code": "file_not_found",
        "message": "File not found"
    }
}
```

### 5. Attempt Path Traversal (Error Path: Invalid Filename)

```bash
curl -v "http://localhost:8085/files/../main.go"
```
**Expected Response (400 Bad Request):**
```json
{
    "error": {
        "code": "invalid_filename",
        "message": "Filename contains invalid characters"
    }
}
```

## Cleanup

-   To remove the dummy file created:
    ```bash
    rm test_upload.txt
    ```
-   To remove the uploaded files:
    ```bash
    rm -rf uploads/*
    ```

## Code Highlights

-   **`main.go`**:
    -   `ctx.Request().ParseMultipartForm(10 << 20)`: Parses the incoming multipart form data with a 10MB memory limit.
    -   `ctx.Request().FormFile("file")`: Retrieves the uploaded file and its header.
    -   `filepath.Base(handler.Filename)` and `strings.Contains(filename, "..")`: Used for basic filename sanitization to prevent directory traversal vulnerabilities.
    -   `os.Create()` and `io.Copy()`: Standard Go functions for saving the uploaded file to disk.
    -   `http.ServeFile()`: Used to efficiently serve static files from the `uploads/` directory.
    -   Error handling is implemented for various stages of the upload and serve process.

This example is essential for applications that need to handle user-generated content like images, documents, or other files.
