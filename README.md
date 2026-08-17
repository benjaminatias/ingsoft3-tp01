# Gestor de Películas

Aplicación web full-stack para mantener una colección personal de películas.

Permite registrar películas, marcarlas como **pendientes** o **vistas**, asignarles una
puntuación, administrar los géneros, filtrar la colección y consultar estadísticas simples.

El acceso está protegido con **autenticación JWT**: hay que crear una cuenta o iniciar
sesión para usar la aplicación.

Los datos se cargan manualmente: no se utiliza ninguna API externa (TMDB, IMDb, OMDb).

> **Nota sobre SPEC.md**: la especificación original excluía explícitamente la
> autenticación, JWT y el registro de usuarios (§2), limitaba el modelo a dos entidades
> (§7) y pedía verificar «No existe autenticación» (§64). El login se agregó después, a
> pedido, y por lo tanto el proyecto ya no cumple esos tres puntos de SPEC.md. Todo lo
> demás de la especificación se mantiene igual.

---

## Stack

```text
Frontend: React + Vite
Backend: Go + Gin
ORM: GORM
Database: PostgreSQL
Backend tests: Go testing
Frontend tests: Vitest
Containers: Docker + Docker Compose
```

## Arquitectura

Monolito sencillo de tres componentes. Sin Repository Pattern y sin microservicios:
los handlers de Gin utilizan GORM directamente. La autenticación es un middleware de
Gin que valida el token antes de llegar a los handlers privados.

```text
React
  │
 REST
  ▼
Go + Gin
  │
 GORM
  ▼
PostgreSQL
```

## Estructura

```text
.
├── README.md
├── .gitignore
├── .env.example
├── docker-compose.yml
│
├── backend/
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── auth/auth.go            (JWT + bcrypt)
│   │   ├── database/database.go
│   │   ├── models/{genero.go,pelicula.go,usuario.go}
│   │   ├── handlers/{handlers.go,auth.go,peliculas.go,generos.go,estadisticas.go,health.go}
│   │   └── validation/validation.go
│   ├── tests/
│   └── Dockerfile
│
└── frontend/
    ├── src/{api,components,utils,App.jsx,main.jsx,styles.css}
    ├── tests/
    ├── Dockerfile
    └── nginx.conf
```

---

## Ejecución con Docker

Linux / macOS:

```bash
cp .env.example .env
docker compose up -d --build
```

Windows (PowerShell):

```powershell
Copy-Item .env.example .env
docker compose up -d --build
```

Windows (CMD):

```cmd
copy .env.example .env
docker compose up -d --build
```

Servicios disponibles:

| Servicio    | URL                                    |
|-------------|----------------------------------------|
| Frontend    | http://localhost:3000                  |
| Backend     | http://localhost:8080                  |
| Películas   | http://localhost:8080/api/peliculas    |
| Géneros     | http://localhost:8080/api/generos      |
| Estadísticas| http://localhost:8080/api/estadisticas |
| Health      | http://localhost:8080/health           |
| PostgreSQL  | localhost:5432                         |

El orden de arranque es: PostgreSQL saludable → backend → backend saludable → frontend.

Detener los servicios conservando los datos:

```bash
docker compose down
```

Detener los servicios y borrar los datos (volumen `db_data`):

```bash
docker compose down -v
```

---

## Ejecución local (sin Docker)

Requiere Go 1.25 o superior, Node 20 o superior y un PostgreSQL accesible con la base
`peliculas_db` creada.

### Backend

```bash
cd backend
go mod download
go run ./cmd/api
```

Variables de entorno utilizadas (con sus valores por defecto entre paréntesis):

```text
DB_HOST       (localhost)
DB_PORT       (5432)
DB_USER       (postgres)
DB_PASSWORD   (vacío)
DB_NAME       (peliculas_db)
DB_SSLMODE    (disable)
SERVER_PORT   (8080)
JWT_SECRET    (si falta, se genera una clave aleatoria solo para desarrollo)
```

`JWT_SECRET` es la clave con la que se firman los tokens. Si no está definida, el backend
genera una aleatoria al arrancar y avisa por consola; al reiniciar, las sesiones anteriores
dejan de ser válidas. Para Docker y para producción hay que definirla siempre
(`openssl rand -base64 32` genera una adecuada).

Las credenciales nunca se hardcodean. Ejemplo en PowerShell:

```powershell
$env:DB_HOST="localhost"; $env:DB_USER="postgres"; $env:DB_PASSWORD="postgres"; $env:DB_NAME="peliculas_db"
go run ./cmd/api
```

Al iniciar, el backend ejecuta las migraciones (`AutoMigrate`) y crea los géneros iniciales.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Disponible en http://localhost:5173. Vite redirige `/api` hacia `http://localhost:8080`,
por eso el frontend siempre usa rutas relativas (`fetch("/api/peliculas")`).

---

## Tests

Backend:

```bash
cd backend
go test ./...
```

Los tests de validación y de handlers no necesitan base de datos. Los tests de
integración (`backend/tests/integracion_test.go`) se ejecutan solamente cuando la
variable `DB_HOST` está definida; en caso contrario se omiten:

```bash
cd backend
DB_HOST=localhost DB_USER=postgres DB_PASSWORD=postgres DB_NAME=peliculas_db go test ./...
```

Frontend:

```bash
cd frontend
npm test -- --run
```

Build:

```bash
cd backend && go build ./cmd/api
cd frontend && npm run build
```

---

## Endpoints

Todas las rutas de la API comienzan con `/api`. Los errores devuelven siempre
`{ "error": "mensaje" }`.

**Salvo el registro, el login y `/health`, todas las rutas exigen el encabezado**
`Authorization: Bearer <token>`. Sin token o con un token vencido devuelven `401`.

### Autenticación (rutas públicas)

| Método | Ruta                  | Descripción                                              |
|--------|-----------------------|----------------------------------------------------------|
| POST   | `/api/auth/registro`  | Crea la cuenta y devuelve el token (201; 409 si el email ya existe) |
| POST   | `/api/auth/login`     | Devuelve el token (401 si las credenciales no son válidas) |
| GET    | `/api/auth/perfil`    | Datos de la cuenta del token (privada)                    |

```bash
# Crear cuenta
curl -X POST http://localhost:8080/api/auth/registro \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Benja","email":"benja@ejemplo.com","password":"contraseña-segura"}'

# Iniciar sesión y usar el token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"benja@ejemplo.com","password":"contraseña-segura"}' | jq -r .token)

curl http://localhost:8080/api/peliculas -H "Authorization: Bearer $TOKEN"
```

Respuesta del registro y del login:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expira": "2026-08-13T20:50:49Z",
  "usuario": { "id": 1, "nombre": "Benja", "email": "benja@ejemplo.com" }
}
```

Reglas de las cuentas:

- Nombre de 2 a 50 caracteres.
- Email válido, único y normalizado a minúsculas (máximo 120 caracteres).
- Contraseña de 8 caracteres como mínimo y 72 bytes como máximo (límite de bcrypt).
- La contraseña se guarda hasheada con bcrypt y nunca se devuelve en las respuestas.
- El token dura 24 horas, se firma con HS256 y solo se acepta ese algoritmo.
- Email inexistente y contraseña incorrecta devuelven el mismo `401` y tardan lo mismo,
  para no revelar qué emails están registrados.

### Películas

| Método | Ruta                             | Descripción                                    |
|--------|----------------------------------|------------------------------------------------|
| GET    | `/api/peliculas`                 | Lista las películas con su género, ordenadas por título |
| GET    | `/api/peliculas/:id`             | Obtiene una película (404 si no existe)        |
| POST   | `/api/peliculas`                 | Crea una película (201)                        |
| PUT    | `/api/peliculas/:id`             | Modifica título, año, género, estado y puntuación |
| DELETE | `/api/peliculas/:id`             | Elimina una película (204)                     |
| PATCH  | `/api/peliculas/:id/estado`      | Cambia el estado; al pasar a pendiente borra la puntuación |
| PATCH  | `/api/peliculas/:id/puntuacion`  | Puntúa la película (400 si está pendiente)     |

Filtros opcionales y combinables de `GET /api/peliculas`:

```text
?estado=vista | pendiente
?generoId=4
?anio=2014
?titulo=inter          (búsqueda parcial, sin distinguir mayúsculas)
?puntuacionMin=8
```

Ejemplo combinado:

```text
GET /api/peliculas?estado=vista&generoId=4&puntuacionMin=8
```

Ejemplo de respuesta:

```json
[
  {
    "id": 1,
    "titulo": "Interstellar",
    "anio": 2014,
    "generoId": 4,
    "genero": { "id": 4, "nombre": "Ciencia ficción" },
    "estado": "vista",
    "puntuacion": 9.5
  }
]
```

### Géneros

| Método | Ruta                 | Descripción                                          |
|--------|----------------------|------------------------------------------------------|
| GET    | `/api/generos`       | Lista los géneros                                     |
| GET    | `/api/generos/:id`   | Obtiene un género (404 si no existe)                  |
| POST   | `/api/generos`       | Crea un género (409 si el nombre ya existe)           |
| PUT    | `/api/generos/:id`   | Edita un género                                       |
| DELETE | `/api/generos/:id`   | Elimina un género (409 si hay películas que lo usan)  |

### Estadísticas y healthcheck

| Método | Ruta                 | Descripción                                          |
|--------|----------------------|------------------------------------------------------|
| GET    | `/api/estadisticas`  | Total, vistas, pendientes y puntuación promedio       |
| GET    | `/health`            | Público (lo usa Docker): 200 si la API y PostgreSQL responden, 503 si no |

```json
{
  "totalPeliculas": 25,
  "vistas": 18,
  "pendientes": 7,
  "puntuacionPromedio": 8.2
}
```

`puntuacionPromedio` considera solamente las películas vistas con puntuación y vale
`null` cuando todavía no hay ninguna.

---

## Reglas de negocio

- Estados permitidos: `pendiente` y `vista`.
- Una película pendiente siempre tiene `puntuacion = null`.
- Solamente se puede puntuar una película vista.
- Al pasar de `vista` a `pendiente` la puntuación se elimina automáticamente.
- Puntuación entre 1 y 10, con un decimal como máximo.
- Título obligatorio, de 1 a 200 caracteres, sin espacios sobrantes.
- Año entre 1888 y el año actual + 5 (el año actual se obtiene con `time.Now().Year()`).
- Género obligatorio y existente.
- Nombres de género de 2 a 50 caracteres, sin duplicados (ignorando mayúsculas y minúsculas).

Géneros creados automáticamente: Acción, Animación, Aventura, Ciencia ficción, Comedia,
Crimen, Documental, Drama, Fantasía, Terror, Thriller y Otros.

La colección es única y compartida por todas las cuentas: el login controla el acceso,
pero no separa la colección por usuario ni existen roles.

---

## Notas

- El archivo `.env` no se sube al repositorio; sí se sube `.env.example`.
- El proyecto está preparado para agregar más adelante `.github/workflows/ci.yml`,
  pero todavía no incluye CI/CD.
