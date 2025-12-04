package templating

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Engine represents the templating engine.
type Engine struct {
	opts  Options
	mu    sync.RWMutex
	cache map[string]*template.Template
}

// NewEngine creates a new templating engine with the given options.
func NewEngine(opts Options) (*Engine, error) {
	if len(opts.Extensions) == 0 {
		opts.Extensions = []string{".gb.html", ".html"}
	}
	return &Engine{
		opts:  opts,
		cache: make(map[string]*template.Template),
	}, nil
}

// Render executes a template with the given name and data, writing the output to w.
func (e *Engine) Render(w io.Writer, name string, data any) error {
	path, err := e.resolvePath(name)
	if err != nil {
		return fmt.Errorf("templating: failed to resolve template path for %s: %w", name, err)
	}

	var tmpl *template.Template
	if e.opts.CacheEnabled && !e.opts.Debug {
		e.mu.RLock()
		tmpl = e.cache[path]
		e.mu.RUnlock()
	}

	if tmpl == nil || e.opts.Debug {
		var loadErr error
		tmpl, loadErr = e.loadTemplate(path)
		if loadErr != nil {
			return fmt.Errorf("templating: failed to load template %s: %w", path, loadErr)
		}
		if e.opts.CacheEnabled && !e.opts.Debug {
			e.mu.Lock()
			e.cache[path] = tmpl
			e.mu.Unlock()
		}
	}

	return tmpl.Execute(w, data)
}

// resolvePath resolves the full path to a template file.
func (e *Engine) resolvePath(name string) (string, error) {
	// If name already ends with one of opts.Extensions, use it as is.
	for _, ext := range e.opts.Extensions {
		if filepath.Ext(name) == ext {
			fullPath := filepath.Join(e.opts.Root, name)
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath, nil
			}
		}
	}

	// Otherwise try appending each extension until a file is found.
	for _, ext := range e.opts.Extensions {
		fullPath := filepath.Join(e.opts.Root, name+ext)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("template %s not found with any of the extensions %v in root %s", name, e.opts.Extensions, e.opts.Root)
}

// loadTemplate reads and parses a template file.
func (e *Engine) loadTemplate(path string) (*template.Template, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file %s: %w", path, err)
	}

	tmpl, err := template.New(filepath.Base(path)).Funcs(e.opts.Funcs).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template file %s: %w", path, err)
	}
	return tmpl, nil
}
