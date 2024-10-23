package main

import (
	"fmt"
	"log"
	"net/http"
	"nms/internal/handlers"
	"nms/internal/handlers/routes"
	"nms/mock"

	"github.com/gorilla/mux"
)

func main() {
	loadDB := mock.LoadMockData()
	fmt.Println(loadDB)
	handler := handlers.NewHandler(loadDB)
	r := mux.NewRouter()

	router := routes.NewHandler(handler, r)

	log.Fatal(http.ListenAndServe(":8080", router))
}
