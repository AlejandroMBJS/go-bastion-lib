package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/response"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

// User represents a simple user model
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// In-memory store for users (for simplicity)
var users = map[int]User{
	1: {ID: 1, Name: "Alice", Age: 30},
	2: {ID: 2, Name: "Bob", Age: 24},
}
var nextUserID = 3

func main() {
	cfg := bastion.DefaultConfig()
	cfg.Port = 8088 // Use a different port for this example

	app := bastion.NewApp(cfg)
	r := app.Router()

	// Global middlewares
	app.Use(
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
	)

	// --- API Endpoints to be tested ---
	r.GET("/users", ListUsersHandler)
	r.POST("/users", CreateUserHandler)
	r.GET("/users/:id", GetUserHandler)
	r.GET("/ping", PingHandler)

	log.Printf("Testing example server starting on :%d", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// PingHandler is a simple handler for /ping
func PingHandler(ctx *router.Context) {
	ctx.JSON(http.StatusOK, map[string]string{"message": "pong"})
}

// ListUsersHandler lists all users
func ListUsersHandler(ctx *router.Context) {
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	response.Success(ctx, http.StatusOK, userList)
}

// CreateUserHandler creates a new user
func CreateUserHandler(ctx *router.Context) {
	var newUser User
	if err := ctx.BindJSON(&newUser); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	if newUser.Name == "" || newUser.Age <= 0 {
		response.Error(ctx, http.StatusBadRequest, "validation_error", "Name and age must be provided and valid")
		return
	}

	newUser.ID = nextUserID
	users[nextUserID] = newUser
	nextUserID++

	response.Created(ctx, fmt.Sprintf("/users/%d", newUser.ID), newUser)
}

// GetUserHandler gets a user by ID
func GetUserHandler(ctx *router.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_id", "Invalid user ID")
		return
	}

	user, ok := users[id]
	if !ok {
		response.Error(ctx, http.StatusNotFound, "not_found", "User not found")
		return
	}
	response.Success(ctx, http.StatusOK, user)
}
