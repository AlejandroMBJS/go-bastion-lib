package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alejandrombjs/go-bastion-lib/pkg/response"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
	"github.com/alejandrombjs/go-bastion-lib/pkg/templating"
)

// MockResponseWriter is a mock implementation of http.ResponseWriter for testing.
type MockResponseWriter struct {
	HeaderMap http.Header
	Body      *bytes.Buffer
	Status    int
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
	}
}

func (m *MockResponseWriter) Header() http.Header {
	return m.HeaderMap
}

func (m *MockResponseWriter) Write(b []byte) (int, error) {
	return m.Body.Write(b)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Status = statusCode
}

func TestHTMLResponse(t *testing.T) {
	// 1. Initialize a default engine with an inline template
	// For testing, we can directly set up the templating engine with a known root
	// and a simple template file.
	testTemplateRoot := t.TempDir()
	templateContent := `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>{{ .Title }}</title>
</head>
<body>
  <h1>{{ .Title }}</h1>
  <p>Hello, {{ .User }}!</p>
</body>
</html>`
	templatePath := fmt.Sprintf("%s/test_home.html", testTemplateRoot)
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test template: %v", err)
	}

	err = templating.InitDefault(templating.Options{
		Root:         testTemplateRoot,
		Extensions:   []string{".gb.html", ".html"}, // Corrected: use double quotes for string literals
		CacheEnabled: false, // Disable cache for testing
		Debug:        true,  // Enable debug for testing
		Funcs:        nil,
	})
	if err != nil {
		t.Fatalf("failed to initialize templating engine: %v", err)
	}

	// 2. Create a bytes.Buffer as a fake http.ResponseWriter (or use httptest.NewRecorder)
	recorder := httptest.NewRecorder()

	// 3. Create a small fake router.Context that uses this recorder.
	// Use router.NewContext to properly initialize the context.
	mockReq := httptest.NewRequest(http.MethodGet, "/", nil)
	mockCtx := router.NewContext(recorder, mockReq)

	// 4. Call response.HTML
	response.HTML(mockCtx, http.StatusOK, "test_home.html", response.H{
		"Title": "Test Page",
		"User":  "Tester",
	})

	// 5. Assertions
	// Status code is 200.
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Content-Type header is text/html; charset=utf-8.
	expectedContentType := "text/html; charset=utf-8"
	if contentType := recorder.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected Content-Type header %q, got %q", expectedContentType, contentType)
	}

	// Body contains <h1>Test Page</h1> and <p>Hello, Tester!</p>.
	expectedBody := `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Test Page</title>
</head>
<body>
  <h1>Test Page</h1>
  <p>Hello, Tester!</p>
</body>
</html>`
	if body := recorder.Body.String(); body != expectedBody {
		t.Errorf("Expected body:\n%q\nGot:\n%q", expectedBody, body)
	}
}
