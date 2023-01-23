package main

import (
	"fmt"
	"net/http"

	service "github.com/sean0427/micro-service-pratice-auth-domain"
	handler "github.com/sean0427/micro-service-pratice-auth-domain/net"
)

func startServer() {
	fmt.Println("Starting server...")

	s := service.New(nil)
	h := handler.New(s)

	handler := h.InitHandler()
	http.ListenAndServe(":8080", handler)

	fmt.Println("Stoping server...")
}

func main() {
	startServer()
}
