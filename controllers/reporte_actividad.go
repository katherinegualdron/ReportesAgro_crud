package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"REPORTES/config"
	"REPORTES/models"

	"github.com/lib/pq"
)

func ObtenerReportesActividad(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`SELECT id_reporte_actividad, tipo_reporte, fecha_inicio, fecha_fin, id_usuario_solicitante, resultado, fecha_generacion FROM "Reportes"."ReporteActividad" ORDER BY id_reporte_actividad`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error al consultar reporte_actividad")
		return
	}
	defer rows.Close()
	items := []models.ReporteActividad{}
	for rows.Next() {
		item, err := scanReporteActividad(rows.Scan)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error al leer reporte_actividad")
			return
		}
		items = append(items, item)
	}
	writeJSON(w, http.StatusOK, items)
}

func ObtenerReporteActividadPorID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "id invalido")
		return
	}
	row := config.DB.QueryRow(`SELECT id_reporte_actividad, tipo_reporte, fecha_inicio, fecha_fin, id_usuario_solicitante, resultado, fecha_generacion FROM "Reportes"."ReporteActividad" WHERE id_reporte_actividad = $1`, id)
	item, err := scanReporteActividad(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "reporte_actividad no encontrado")
			return
		}
		writeError(w, http.StatusInternalServerError, "error al consultar reporte_actividad")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func CrearReporteActividad(w http.ResponseWriter, r *http.Request) {
	var item models.ReporteActividad
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "json invalido")
		return
	}
	if ok, err := existeRegistro(`"Usuarios"."Usuario"`, "id_usuario", item.IdUsuarioSolicitante); err != nil {
		writeError(w, http.StatusBadRequest, "error al validar id_usuario_solicitante")
		return
	} else if !ok {
		writeError(w, http.StatusBadRequest, "id_usuario_solicitante no existe")
		return
	}
	row := config.DB.QueryRow(`INSERT INTO "Reportes"."ReporteActividad" (tipo_reporte, fecha_inicio, fecha_fin, id_usuario_solicitante, resultado, fecha_generacion) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_reporte_actividad, tipo_reporte, fecha_inicio, fecha_fin, id_usuario_solicitante, resultado, fecha_generacion`,
		item.TipoReporte, item.FechaInicio, item.FechaFin, item.IdUsuarioSolicitante, item.Resultado, item.FechaGeneracion)
	item, err := scanReporteActividad(row.Scan)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23514" {
			writeError(w, http.StatusBadRequest, "fecha_fin no puede ser menor que fecha_inicio")
			return
		}
		writeError(w, http.StatusInternalServerError, "error al crear reporte_actividad")
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func ActualizarReporteActividad(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "id invalido")
		return
	}
	var item models.ReporteActividad
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "json invalido")
		return
	}
	if ok, err := existeRegistro(`"Usuarios"."Usuario"`, "id_usuario", item.IdUsuarioSolicitante); err != nil {
		writeError(w, http.StatusBadRequest, "error al validar id_usuario_solicitante")
		return
	} else if !ok {
		writeError(w, http.StatusBadRequest, "id_usuario_solicitante no existe")
		return
	}
	row := config.DB.QueryRow(`UPDATE "Reportes"."ReporteActividad" SET tipo_reporte = $1, fecha_inicio = $2, fecha_fin = $3, id_usuario_solicitante = $4, resultado = $5, fecha_generacion = $6 WHERE id_reporte_actividad = $7 RETURNING id_reporte_actividad, tipo_reporte, fecha_inicio, fecha_fin, id_usuario_solicitante, resultado, fecha_generacion`,
		item.TipoReporte, item.FechaInicio, item.FechaFin, item.IdUsuarioSolicitante, item.Resultado, item.FechaGeneracion, id)
	item, err = scanReporteActividad(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "reporte_actividad no encontrado")
			return
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23514" {
			writeError(w, http.StatusBadRequest, "fecha_fin no puede ser menor que fecha_inicio")
			return
		}
		writeError(w, http.StatusInternalServerError, "error al actualizar reporte_actividad")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func EliminarReporteActividad(w http.ResponseWriter, r *http.Request) {
	eliminarGenericoReportes(w, r, `"Reportes"."ReporteActividad"`, "id_reporte_actividad", "reporte_actividad")
}

func scanReporteActividad(scan func(dest ...any) error) (models.ReporteActividad, error) {
	var item models.ReporteActividad
	var resultado sql.NullString
	err := scan(&item.IdReporteActividad, &item.TipoReporte, &item.FechaInicio, &item.FechaFin, &item.IdUsuarioSolicitante, &resultado, &item.FechaGeneracion)
	if err != nil {
		return models.ReporteActividad{}, err
	}
	item.Resultado = nullStringToPointer(resultado)
	return item, nil
}
