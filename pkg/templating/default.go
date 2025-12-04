package templating

import (
	"fmt"
	"io"
	"sync"
)

var (
	defaultEngine *Engine
	defaultMu     sync.RWMutex
)

// InitDefault initializes the default templating engine.
func InitDefault(opts Options) error {
	defaultMu.Lock()
	defer defaultMu.Unlock()

	eng, err := NewEngine(opts)
	if err != nil {
		return fmt.Errorf("templating: failed to create default engine: %w", err)
	}
	defaultEngine = eng
	return nil
}

// MustInitDefault initializes the default templating engine and panics on error.
func MustInitDefault(opts Options) {
	if err := InitDefault(opts); err != nil {
		panic(err)
	}
}

// getDefault returns the default templating engine or an error if not initialized.
func getDefault() (*Engine, error) {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	if defaultEngine == nil {
		return nil, fmt.Errorf("templating: default engine not initialized")
	}
	return defaultEngine, nil
}

// Render renders a template using the default engine.
func Render(w io.Writer, name string, data any) error {
	eng, err := getDefault()
	if err != nil {
		return err
	}
	return eng.Render(w, name, data)
}
