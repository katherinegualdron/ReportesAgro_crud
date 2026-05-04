package main

import (
	"log"
	"net/http"

	"REPORTES/config"
	"REPORTES/routes"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	router := mux.NewRouter()
	routes.RegistrarRutas(router)

	log.Println("Servidor REPORTES escuchando en :8093")
	log.Fatal(http.ListenAndServe(":8093", router))
}
