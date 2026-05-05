# ReportesAgro CRUD

API REST desarrollada en Go para gestionar reportes de actividad y exportaciones de datos del sistema Agrocampo.

El proyecto expone endpoints CRUD para dos recursos principales:

- `reporte_actividad`: administra reportes generados por usuarios.
- `tr_exportacion_datos`: administra exportaciones asociadas a reportes.

## Tecnologias

- Go 1.22
- Gorilla Mux
- PostgreSQL
- Driver `github.com/lib/pq`

## Estructura del proyecto

```text
.
|-- config/
|   `-- db.go
|-- controllers/
|   |-- helpers.go
|   |-- reporte_actividad.go
|   `-- tr_exportacion_datos.go
|-- models/
|   |-- reporte_actividad.go
|   `-- tr_exportacion_datos.go
|-- routes/
|   |-- routes.go
|   |-- reporte_actividad.go
|   `-- tr_exportacion_datos.go
|-- go.mod
`-- main.go
```

## Conexion a base de datos

La conexion se configura en `config/db.go`.

Valores actuales:

```text
host: localhost
port: 5432
user: postgres
password: postgres
database: Agrocampo
schema: Reportes
```

El codigo consulta las tablas con el esquema y nombre exacto:

- `"Reportes"."ReporteActividad"`
- `"Reportes"."tr_ExportacionDatos"`
- `"Usuarios"."Usuario"`

Esto es importante porque PostgreSQL diferencia mayusculas y minusculas cuando los nombres estan entre comillas.

## Ejecutar el proyecto

Instalar dependencias:

```bash
go mod tidy
```

Ejecutar el servidor:

```bash
go run main.go
```

El servidor queda escuchando en:

```text
http://localhost:8095
```

Para detenerlo, usar `Ctrl + C` en la terminal.

## Endpoints

### Reporte actividad

| Metodo | Ruta | Descripcion |
| --- | --- | --- |
| GET | `/reporte_actividad` | Lista todos los reportes de actividad |
| GET | `/reporte_actividad/{id}` | Obtiene un reporte por ID |
| POST | `/reporte_actividad` | Crea un reporte |
| PUT | `/reporte_actividad/{id}` | Actualiza un reporte |
| DELETE | `/reporte_actividad/{id}` | Elimina un reporte |

Ejemplo para crear o actualizar:

```json
{
  "tipo_reporte": "actividad",
  "fecha_inicio": "2026-05-01T00:00:00Z",
  "fecha_fin": "2026-05-04T00:00:00Z",
  "id_usuario_solicitante": 1,
  "resultado": "Reporte generado correctamente",
  "fecha_generacion": "2026-05-04T19:00:00Z"
}
```

Validaciones principales:

- `id_usuario_solicitante` debe existir en `"Usuarios"."Usuario"`.
- `fecha_fin` no puede ser menor que `fecha_inicio`.
- `resultado` puede enviarse como texto o `null`.

### Exportacion de datos

| Metodo | Ruta | Descripcion |
| --- | --- | --- |
| GET | `/tr_exportacion_datos` | Lista todas las exportaciones |
| GET | `/tr_exportacion_datos/{id}` | Obtiene una exportacion por ID |
| POST | `/tr_exportacion_datos` | Crea una exportacion |
| PUT | `/tr_exportacion_datos/{id}` | Actualiza una exportacion |
| DELETE | `/tr_exportacion_datos/{id}` | Elimina una exportacion |

Ejemplo para crear o actualizar:

```json
{
  "id_usuario": 1,
  "id_reporte": 1,
  "formato": "PDF",
  "url_archivo": "https://ejemplo.com/reportes/reporte-1.pdf",
  "estado": "GENERADO",
  "fecha_creacion": "2026-05-04T19:00:00Z"
}
```

Validaciones principales:

- `id_usuario` debe existir en `"Usuarios"."Usuario"`.
- `id_reporte` debe existir en `"Reportes"."ReporteActividad"`.
- `formato` y `estado` deben cumplir las restricciones definidas en la base de datos.
- `url_archivo` puede enviarse como texto o `null`.

## Respuestas

Respuesta exitosa:

```json
{
  "id_reporte_actividad": 1,
  "tipo_reporte": "actividad",
  "fecha_inicio": "2026-05-01T00:00:00Z",
  "fecha_fin": "2026-05-04T00:00:00Z",
  "id_usuario_solicitante": 1,
  "resultado": "Reporte generado correctamente",
  "fecha_generacion": "2026-05-04T19:00:00Z"
}
```

Respuesta de error:

```json
{
  "error": "mensaje del error"
}
```

## Codigos HTTP usados

| Codigo | Uso |
| --- | --- |
| 200 | Consulta, actualizacion o eliminacion exitosa |
| 201 | Registro creado correctamente |
| 400 | JSON invalido, ID invalido o validacion fallida |
| 404 | Registro no encontrado |
| 500 | Error interno al consultar, crear, actualizar o eliminar |

## Pruebas con Postman

Usar como base:

```text
http://localhost:8095
```

Ejemplo:

```text
GET http://localhost:8095/reporte_actividad
```

Para solicitudes `POST` y `PUT`, seleccionar:

```text
Body -> raw -> JSON
```

Tambien agregar el header:

```text
Content-Type: application/json
```

## Errores comunes

Si aparece un error parecido a:

```text
listen tcp :8095: bind: Only one usage of each socket address is normally permitted
```

significa que el servidor ya esta corriendo en ese puerto. Cerrar la terminal anterior con `Ctrl + C` o finalizar el proceso que usa el puerto.

Si aparece:

```json
{
  "error": "error al consultar reporte_actividad"
}
```

revisar que la base de datos `Agrocampo` exista, que PostgreSQL este encendido y que las tablas existan con estos nombres exactos:

- `"Reportes"."ReporteActividad"`
- `"Reportes"."tr_ExportacionDatos"`

## Verificacion

Para comprobar que el proyecto compila:

```bash
go test ./...
```
