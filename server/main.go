package main

import (
	"github.com/KouT127/pub-sub-practice/config"
	"github.com/KouT127/pub-sub-practice/handler"
	"github.com/KouT127/pub-sub-practice/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	config.Configure()
	http.HandleFunc("/resize", middleware.CloudTraceMiddleware(handler.ResizeSubscriptionHandler))

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
