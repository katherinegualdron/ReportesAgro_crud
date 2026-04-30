package routes

import (
	"REPORTES/controllers"

	"github.com/gorilla/mux"
)

func RegistrarRutasReporteActividad(router *mux.Router) {
	router.HandleFunc("/reporte_actividad", controllers.ObtenerReportesActividad).Methods("GET")
	router.HandleFunc("/reporte_actividad/{id}", controllers.ObtenerReporteActividadPorID).Methods("GET")
	router.HandleFunc("/reporte_actividad", controllers.CrearReporteActividad).Methods("POST")
	router.HandleFunc("/reporte_actividad/{id}", controllers.ActualizarReporteActividad).Methods("PUT")
	router.HandleFunc("/reporte_actividad/{id}", controllers.EliminarReporteActividad).Methods("DELETE")
}
