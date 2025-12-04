package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/response"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

const uploadDir = "./uploads" // Directory to store uploaded files

func main() {
	cfg := bastion.DefaultConfig()
	cfg.Port = 8085 // Use a different port for this example

	app := bastion.NewApp(cfg)
	r := app.Router()

	// Global middlewares
	app.Use(
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
	)

	// Ensure upload directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.Mkdir(uploadDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create upload directory: %v", err)
		}
	}

	// --- File Upload Endpoint ---
	r.POST("/upload", func(ctx *router.Context) {
		// Parse multipart form, 10MB limit
		err := ctx.Request().ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid_form", "Failed to parse multipart form: "+err.Error())
			return
		}

		// Get the file from the form data
		file, handler, err := ctx.Request().FormFile("file")
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "missing_file", "Error retrieving file from form: "+err.Error())
			return
		}
		defer file.Close()

		// Sanitize filename to prevent path traversal attacks
		filename := filepath.Base(handler.Filename)
		if strings.Contains(filename, "..") {
			response.Error(ctx, http.StatusBadRequest, "invalid_filename", "Filename contains invalid characters")
			return
		}

		// Create a new file in the uploads directory
		dstPath := filepath.Join(uploadDir, filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			response.Error(ctx, http.StatusInternalServerError, "file_create_error", "Failed to create file on server: "+err.Error())
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the destination
		if _, err := io.Copy(dst, file); err != nil {
			response.Error(ctx, http.StatusInternalServerError, "file_copy_error", "Failed to save file on server: "+err.Error())
			return
		}

		response.Success(ctx, http.StatusOK, map[string]string{
			"message":  "File uploaded successfully",
			"filename": filename,
			"size":     fmt.Sprintf("%d bytes", handler.Size),
			"location": dstPath,
		})
	})

	// --- Serve Uploaded Files Endpoint ---
	r.GET("/files/:name", func(ctx *router.Context) {
		filename := ctx.Param("name")
		if strings.Contains(filename, "..") { // Prevent path traversal
			response.Error(ctx, http.StatusBadRequest, "invalid_filename", "Filename contains invalid characters")
			return
		}

		filePath := filepath.Join(uploadDir, filename)

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			response.Error(ctx, http.StatusNotFound, "file_not_found", "File not found")
			return
		}

		// Serve the file
		http.ServeFile(ctx.ResponseWriter(), ctx.Request().Request, filePath)
	})

	log.Printf("File Upload example server starting on :%d", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
