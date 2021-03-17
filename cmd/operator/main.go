package main

import (
	"log"
	"net/http"

	"github.com/shiftstack/metrics"
)

func main() {

	http.Handle("/metrics", metrics.New())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
