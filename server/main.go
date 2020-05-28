package main

import (
	"github.com/KouT127/cloud-playground/config"
	"github.com/KouT127/cloud-playground/handler"
	"github.com/KouT127/cloud-playground/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	config.Configure()
	http.HandleFunc("/resize", middleware.CloudTraceMiddleware(handler.ResizeSubscriptionHandler))
	http.HandleFunc("/task", middleware.CloudTraceMiddleware(handler.TaskHandler))


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
