package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/uchitsa/ubiquitous-octo-train/internal/app"
)

func main() {
	r := chi.NewRouter()
	orderService := app.NewOrderService()

	r.Post("/orders", orderService.CreateOrder)

	log.Println("Server listening on localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Server failed: %s", err)
		os.Exit(1)
	}
}
