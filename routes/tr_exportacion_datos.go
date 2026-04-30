package routes

import (
	"REPORTES/controllers"

	"github.com/gorilla/mux"
)

func RegistrarRutasExportacionDatos(router *mux.Router) {
	router.HandleFunc("/tr_exportacion_datos", controllers.ObtenerExportacionesDatos).Methods("GET")
	router.HandleFunc("/tr_exportacion_datos/{id}", controllers.ObtenerExportacionDatosPorID).Methods("GET")
	router.HandleFunc("/tr_exportacion_datos", controllers.CrearExportacionDatos).Methods("POST")
	router.HandleFunc("/tr_exportacion_datos/{id}", controllers.ActualizarExportacionDatos).Methods("PUT")
	router.HandleFunc("/tr_exportacion_datos/{id}", controllers.EliminarExportacionDatos).Methods("DELETE")
}
