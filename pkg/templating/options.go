package templating

import (
	"html/template"
)

// Options holds the configuration for the templating engine.
type Options struct {
	Root         string              // templates root directory
	Extensions   []string            // e.g. []string{".gb.html", ".html"}
	Funcs        template.FuncMap    // extra functions
	CacheEnabled bool                // cache parsed templates (prod)
	Debug        bool                // if true, reload templates every time
}

// H is a convenient data type for passing data to templates.
type H map[string]any
