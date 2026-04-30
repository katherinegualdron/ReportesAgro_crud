package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"REPORTES/config"

	"github.com/gorilla/mux"
)

func getIDFromRequest(r *http.Request, key string) (int, error) {
	return strconv.Atoi(mux.Vars(r)[key])
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}

func existeRegistro(tabla string, columna string, id int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM " + tabla + " WHERE " + columna + " = $1)"
	var existe bool
	if err := config.DB.QueryRow(query, id).Scan(&existe); err != nil {
		return false, err
	}
	return existe, nil
}

func nullStringToPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	text := value.String
	return &text
}

func eliminarGenericoReportes(w http.ResponseWriter, r *http.Request, tabla string, columna string, nombre string) {
	id, err := getIDFromRequest(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "id invalido")
		return
	}
	result, err := config.DB.Exec("DELETE FROM "+tabla+" WHERE "+columna+" = $1", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error al eliminar "+nombre)
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, nombre+" no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": nombre + " eliminado correctamente"})
}
