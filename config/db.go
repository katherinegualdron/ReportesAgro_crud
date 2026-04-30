package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "postgres"
	dbname := "AgroCampo"
	schema := "Reportes"

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		host, port, user, password, dbname, schema,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error al conectar:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("no se puede conectar:", err)
	}

	fmt.Println("conexion a base de datos exitosa")
	fmt.Println("Conectado a la db:", dbname, "y esquema:", schema)
	DB = db
}
