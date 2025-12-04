# GORM + PostgreSQL CRUD API Example

This example demonstrates how to build a complete CRUD (Create, Read, Update, Delete) API for a `Product` resource using `goBastion Lib`, GORM (Go Object Relational Mapper), and a PostgreSQL database. All services are orchestrated with Docker Compose.

## Features

-   CRUD operations for a `Product` resource:
    -   `GET /products`: List all products.
    -   `POST /products`: Create a new product.
    -   `GET /products/:id`: Retrieve a single product by ID.
    -   `PUT /products/:id`: Update an existing product by ID.
    -   `DELETE /products/:id`: Delete a product by ID.
-   Database persistence using PostgreSQL.
-   GORM for ORM capabilities and schema migration.
-   Error handling for invalid input, not found resources, and database errors.
-   Docker Compose setup for easy local development.

## Prerequisites

-   Docker and Docker Compose installed.
-   `jq` for parsing JSON responses in test commands.

## How to Run

1.  Navigate to the example directory:
    ```bash
    cd examples/gorm-postgres-crud
    ```
2.  Build and start the services using Docker Compose:
    ```bash
    docker-compose up --build -d
    ```
    This will:
    -   Build the Go API service (`api`) from the `Dockerfile`.
    -   Start a PostgreSQL database service (`db`).
    -   The Go application will connect to the database, perform auto-migration (creating the `products` table), and start on port `8084`.

3.  Verify the services are running:
    ```bash
    docker-compose ps
    ```
    You should see `api` and `db` services in a healthy state.

## How to Test

The API service will be accessible at `http://localhost:8084`.

### 1. Create a Product (POST /products)

```bash
CREATE_PRODUCT_RESPONSE=$(curl -s -X POST "http://localhost:8084/products" \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop Pro","description":"Powerful laptop for professionals","price":1200.00,"stock":50}')

echo "$CREATE_PRODUCT_RESPONSE"
PRODUCT_ID=$(echo "$CREATE_PRODUCT_RESPONSE" | jq -r '.id')
echo "Created Product ID: $PRODUCT_ID"
```
**Expected Response (201 Created):**
```json
{
    "id": 1,
    "name": "Laptop Pro",
    "description": "Powerful laptop for professionals",
    "price": 1200,
    "stock": 50,
    "CreatedAt": "...",
    "UpdatedAt": "...",
    "DeletedAt": null
}
```

### 2. List All Products (GET /products)

```bash
curl -v "http://localhost:8084/products"
```
**Expected Response (200 OK):**
```json
[
    {
        "id": 1,
        "name": "Laptop Pro",
        "description": "Powerful laptop for professionals",
        "price": 1200,
        "stock": 50,
        "CreatedAt": "...",
        "UpdatedAt": "...",
        "DeletedAt": null
    }
]
```

### 3. Get Product by ID (GET /products/:id)

```bash
curl -v "http://localhost:8084/products/$PRODUCT_ID"
```
**Expected Response (200 OK):**
```json
{
    "id": 1,
    "name": "Laptop Pro",
    "description": "Powerful laptop for professionals",
    "price": 1200,
    "stock": 50,
    "CreatedAt": "...",
    "UpdatedAt": "...",
    "DeletedAt": null
}
```

### 4. Update Product (PUT /products/:id)

```bash
curl -v -X PUT "http://localhost:8084/products/$PRODUCT_ID" \
  -H "Content-Type: application/json" \
  -d '{"price":1150.00,"stock":45}'
```
**Expected Response (200 OK):**
```json
{
    "id": 1,
    "name": "Laptop Pro",
    "description": "Powerful laptop for professionals",
    "price": 1150,
    "stock": 45,
    "CreatedAt": "...",
    "UpdatedAt": "...",
    "DeletedAt": null
}
```

### 5. Delete Product (DELETE /products/:id)

```bash
curl -v -X DELETE "http://localhost:8084/products/$PRODUCT_ID"
```
**Expected Response (204 No Content):** (No body returned)

### 6. Get Deleted Product (Error Path: Not Found)

```bash
curl -v "http://localhost:8084/products/$PRODUCT_ID"
```
**Expected Response (404 Not Found):**
```json
{
    "error": {
        "code": "not_found",
        "message": "Product not found"
    }
}
```

### 7. Create Product (Error Path: Validation Error)

```bash
curl -v -X POST "http://localhost:8084/products" \
  -H "Content-Type: application/json" \
  -d '{"name":"","description":"Invalid product","price":-10.00,"stock":5}'
```
**Expected Response (400 Bad Request):**
```json
{
    "error": {
        "code": "validation_error",
        "message": "Product name and price must be provided and valid"
    }
}
```

## Cleanup

To stop and remove the Docker containers and associated volumes:
```bash
docker-compose down -v
```

## Code Highlights

-   **`main.go`**:
    -   Connects to PostgreSQL using GORM and performs `db.AutoMigrate(&Product{})`.
    -   Defines handlers for each CRUD operation (`GET`, `POST`, `PUT`, `DELETE`) on the `/products` endpoint.
    -   Uses `ctx.Param("id")` to extract product IDs from the URL.
    -   Utilizes `ctx.BindJSON()` for request body parsing and `response.Success()`, `response.Created()`, `response.NoContent()`, and `response.Error()` for consistent JSON responses.
    -   Includes retry logic for database connection, useful in Docker Compose environments where the DB might not be immediately ready.
-   **`docker-compose.yml`**: Sets up the `api` (Go app) and `db` (PostgreSQL) services, managing their environment variables and networking.
-   **`Dockerfile`**: Builds the Go application into a lightweight Alpine-based image.

This example provides a solid foundation for building data-driven APIs with `goBastion Lib` and GORM.
