package main

import (
	"fmt"
	"log"
	"net/http"

	"imaging-service/internal/handler"
)

func main() {
	http.HandleFunc("/", handler.ImageHandler)
	port := 8080
	fmt.Printf("Imaging Service running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
