package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/response"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
)

func main() {
	cfg := bastion.DefaultConfig()
	cfg.Port = 9876 // Ensure the port is 9876 for the example
	cfg.TemplateRoot = "examples/html_example/templates" // Set template root for this example

	app := bastion.NewApp(cfg)
	r := app.Router()

	r.GET("/", func(ctx *router.Context) {
		response.HTML(ctx, http.StatusOK, "home.html", response.H{
			"Title": "goBastion HTML Test",
			"User":  "Alejandro",
		})
	})

	log.Printf("HTML example listening on :%d\n", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r.Handler()); err != nil {
		log.Fatal(err)
	}
}
