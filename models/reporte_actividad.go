package models

import "time"

type ReporteActividad struct {
	IdReporteActividad   int       `json:"id_reporte_actividad"`
	TipoReporte          string    `json:"tipo_reporte"`
	FechaInicio          time.Time `json:"fecha_inicio"`
	FechaFin             time.Time `json:"fecha_fin"`
	IdUsuarioSolicitante int       `json:"id_usuario_solicitante"`
	Resultado            *string   `json:"resultado"`
	FechaGeneracion      time.Time `json:"fecha_generacion"`
}
