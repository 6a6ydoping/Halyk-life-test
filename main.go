package main

import (
	"halykTestTask/controllers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", controllers.HandleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
