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

func ObtenerExportacionesDatos(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`SELECT id_tr_exportacion_datos, id_usuario, id_reporte, formato, url_archivo, estado, fecha_creacion FROM "Reportes"."tr_ExportacionDatos" ORDER BY id_tr_exportacion_datos`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error al consultar tr_exportacion_datos")
		return
	}
	defer rows.Close()
	items := []models.TrExportacionDatos{}
	for rows.Next() {
		item, err := scanExportacionDatos(rows.Scan)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error al leer tr_exportacion_datos")
			return
		}
		items = append(items, item)
	}
	writeJSON(w, http.StatusOK, items)
}

func ObtenerExportacionDatosPorID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "id invalido")
		return
	}
	row := config.DB.QueryRow(`SELECT id_tr_exportacion_datos, id_usuario, id_reporte, formato, url_archivo, estado, fecha_creacion FROM "Reportes"."tr_ExportacionDatos" WHERE id_tr_exportacion_datos = $1`, id)
	item, err := scanExportacionDatos(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "tr_exportacion_datos no encontrado")
			return
		}
		writeError(w, http.StatusInternalServerError, "error al consultar tr_exportacion_datos")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func CrearExportacionDatos(w http.ResponseWriter, r *http.Request) {
	var item models.TrExportacionDatos
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "json invalido")
		return
	}
	if err := validarExportacionDatos(item.IdUsuario, item.IdReporte); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	row := config.DB.QueryRow(`INSERT INTO "Reportes"."tr_ExportacionDatos" (id_usuario, id_reporte, formato, url_archivo, estado, fecha_creacion) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_tr_exportacion_datos, id_usuario, id_reporte, formato, url_archivo, estado, fecha_creacion`,
		item.IdUsuario, item.IdReporte, item.Formato, item.UrlArchivo, item.Estado, item.FechaCreacion)
	item, err := scanExportacionDatos(row.Scan)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23514" {
			writeError(w, http.StatusBadRequest, "formato o estado no valido")
			return
		}
		writeError(w, http.StatusInternalServerError, "error al crear tr_exportacion_datos")
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func ActualizarExportacionDatos(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "id invalido")
		return
	}
	var item models.TrExportacionDatos
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "json invalido")
		return
	}
	if err := validarExportacionDatos(item.IdUsuario, item.IdReporte); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	row := config.DB.QueryRow(`UPDATE "Reportes"."tr_ExportacionDatos" SET id_usuario = $1, id_reporte = $2, formato = $3, url_archivo = $4, estado = $5, fecha_creacion = $6 WHERE id_tr_exportacion_datos = $7 RETURNING id_tr_exportacion_datos, id_usuario, id_reporte, formato, url_archivo, estado, fecha_creacion`,
		item.IdUsuario, item.IdReporte, item.Formato, item.UrlArchivo, item.Estado, item.FechaCreacion, id)
	item, err = scanExportacionDatos(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "tr_exportacion_datos no encontrado")
			return
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23514" {
			writeError(w, http.StatusBadRequest, "formato o estado no valido")
			return
		}
		writeError(w, http.StatusInternalServerError, "error al actualizar tr_exportacion_datos")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func EliminarExportacionDatos(w http.ResponseWriter, r *http.Request) {
	eliminarGenericoReportes(w, r, `"Reportes"."tr_ExportacionDatos"`, "id_tr_exportacion_datos", "tr_exportacion_datos")
}

func scanExportacionDatos(scan func(dest ...any) error) (models.TrExportacionDatos, error) {
	var item models.TrExportacionDatos
	var urlArchivo sql.NullString
	err := scan(&item.IdTrExportacionDatos, &item.IdUsuario, &item.IdReporte, &item.Formato, &urlArchivo, &item.Estado, &item.FechaCreacion)
	if err != nil {
		return models.TrExportacionDatos{}, err
	}
	item.UrlArchivo = nullStringToPointer(urlArchivo)
	return item, nil
}

func validarExportacionDatos(idUsuario int, idReporte int) error {
	if ok, err := existeRegistro(`"Usuarios"."Usuario"`, "id_usuario", idUsuario); err != nil {
		return errors.New("error al validar id_usuario")
	} else if !ok {
		return errors.New("id_usuario no existe")
	}
	if ok, err := existeRegistro(`"Reportes"."ReporteActividad"`, "id_reporte_actividad", idReporte); err != nil {
		return errors.New("error al validar id_reporte")
	} else if !ok {
		return errors.New("id_reporte no existe")
	}
	return nil
}
