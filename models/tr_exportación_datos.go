package models

import "time"

type TrExportacionDatos struct {
	IdTrExportacionDatos int       `json:"id_tr_exportacion_datos"`
	IdUsuario            int       `json:"id_usuario"`
	IdReporte            int       `json:"id_reporte"`
	Formato              string    `json:"formato"`
	UrlArchivo           *string   `json:"url_archivo"`
	Estado               string    `json:"estado"`
	FechaCreacion        time.Time `json:"fecha_creacion"`
}
