package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gestor-peliculas/internal/database"
	"gestor-peliculas/internal/handlers"
	"gestor-peliculas/internal/models"
)

const (
	emailDePrueba    = "tests-integracion@ejemplo.com"
	passwordDePrueba = "contraseña-de-prueba"
)

// prepararIntegracion devuelve un router conectado a PostgreSQL y un token válido.
// Estos tests solamente se ejecutan cuando existe la variable DB_HOST,
// de manera que `go test ./...` funcione también sin base de datos.
func prepararIntegracion(t *testing.T) (*gin.Engine, *gorm.DB, string) {
	t.Helper()

	if _, definida := os.LookupEnv("DB_HOST"); !definida {
		t.Skip("DB_HOST no está definida: se omiten los tests de integración")
	}

	db, err := database.Conectar()
	if err != nil {
		t.Skipf("no hay PostgreSQL disponible: %v", err)
	}
	if err := database.Migrar(db); err != nil {
		t.Fatalf("no se pudieron ejecutar las migraciones: %v", err)
	}
	if err := database.CrearGenerosIniciales(db); err != nil {
		t.Fatalf("no se pudieron crear los géneros iniciales: %v", err)
	}

	router := handlers.NuevoRouter(db)

	// La cuenta de pruebas se crea desde cero en cada corrida.
	db.Unscoped().Where("email = ?", emailDePrueba).Delete(&models.Usuario{})

	cuerpo := `{"nombre":"Tester","email":"` + emailDePrueba + `","password":"` + passwordDePrueba + `"}`
	respuesta := ejecutar(t, router, http.MethodPost, "/api/auth/registro", cuerpo)
	if respuesta.Code != http.StatusCreated {
		t.Fatalf("no se pudo registrar la cuenta de pruebas: %d (%s)", respuesta.Code, respuesta.Body.String())
	}

	var sesion struct {
		Token   string         `json:"token"`
		Usuario models.Usuario `json:"usuario"`
	}
	decodificar(t, respuesta.Body.Bytes(), &sesion)
	if sesion.Token == "" {
		t.Fatal("el registro debe devolver un token")
	}
	if sesion.Usuario.PasswordHash != "" {
		t.Fatal("la respuesta nunca debe incluir el hash de la contraseña")
	}

	t.Cleanup(func() {
		db.Unscoped().Where("email = ?", emailDePrueba).Delete(&models.Usuario{})
	})

	return router, db, sesion.Token
}

// conSesion lanza una petición autenticada con el token de la sesión de pruebas.
func conSesion(t *testing.T, router *gin.Engine, token, metodo, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	return ejecutarConCabecera(t, router, metodo, ruta, cuerpo, "Bearer "+token)
}

func decodificar(t *testing.T, cuerpo []byte, destino any) {
	t.Helper()
	if err := json.Unmarshal(cuerpo, destino); err != nil {
		t.Fatalf("la respuesta no es JSON válido: %v", err)
	}
}

func TestIntegracionCicloCompletoDePelicula(t *testing.T) {
	router, db, token := prepararIntegracion(t)

	var genero models.Genero
	if err := db.Where("nombre = ?", "Ciencia ficción").First(&genero).Error; err != nil {
		t.Fatalf("no se encontró el género inicial: %v", err)
	}

	// Crear una película vista con puntuación.
	cuerpo := `{"titulo":"  Pelicula De Prueba  ","anio":2014,"generoId":` +
		itoa(genero.ID) + `,"estado":"vista","puntuacion":9.5}`
	respuesta := conSesion(t, router, token, http.MethodPost, "/api/peliculas", cuerpo)
	if respuesta.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
	}

	var creada models.Pelicula
	decodificar(t, respuesta.Body.Bytes(), &creada)
	defer db.Unscoped().Delete(&models.Pelicula{}, creada.ID)

	if creada.Titulo != "Pelicula De Prueba" {
		t.Fatalf("el título debería quedar limpio, se obtuvo %q", creada.Titulo)
	}
	if creada.Genero.Nombre != "Ciencia ficción" {
		t.Fatalf("la respuesta debe incluir el género, se obtuvo %q", creada.Genero.Nombre)
	}
	if creada.Puntuacion == nil || *creada.Puntuacion != 9.5 {
		t.Fatal("la puntuación no se guardó correctamente")
	}

	// Obtenerla por identificador.
	respuesta = conSesion(t, router, token, http.MethodGet, "/api/peliculas/"+itoa(creada.ID), "")
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", respuesta.Code)
	}

	// Buscarla por título de forma parcial e insensible a mayúsculas.
	respuesta = conSesion(t, router, token, http.MethodGet, "/api/peliculas?titulo=de+prueba&estado=vista&puntuacionMin=9", "")
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", respuesta.Code)
	}
	var encontradas []models.Pelicula
	decodificar(t, respuesta.Body.Bytes(), &encontradas)
	if len(encontradas) == 0 {
		t.Fatal("la búsqueda parcial por título debería encontrar la película")
	}

	// Pasar a pendiente: la puntuación debe eliminarse automáticamente.
	respuesta = conSesion(t, router, token, http.MethodPatch, "/api/peliculas/"+itoa(creada.ID)+"/estado", `{"estado":"pendiente"}`)
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
	}
	var actualizada models.Pelicula
	decodificar(t, respuesta.Body.Bytes(), &actualizada)
	if actualizada.Estado != "pendiente" {
		t.Fatalf("se esperaba estado pendiente, se obtuvo %q", actualizada.Estado)
	}
	if actualizada.Puntuacion != nil {
		t.Fatalf("al pasar a pendiente la puntuación debe quedar en null, se obtuvo %v", *actualizada.Puntuacion)
	}

	// Puntuar una película pendiente no está permitido.
	respuesta = conSesion(t, router, token, http.MethodPatch, "/api/peliculas/"+itoa(creada.ID)+"/puntuacion", `{"puntuacion":8.5}`)
	if respuesta.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400 al puntuar una película pendiente, se obtuvo %d", respuesta.Code)
	}

	// Volver a vista y puntuar.
	respuesta = conSesion(t, router, token, http.MethodPatch, "/api/peliculas/"+itoa(creada.ID)+"/estado", `{"estado":"vista"}`)
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", respuesta.Code)
	}
	respuesta = conSesion(t, router, token, http.MethodPatch, "/api/peliculas/"+itoa(creada.ID)+"/puntuacion", `{"puntuacion":8.5}`)
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
	}

	// Modificar la película completa, incluyendo el género.
	var otroGenero models.Genero
	if err := db.Where("nombre = ?", "Drama").First(&otroGenero).Error; err != nil {
		t.Fatalf("no se encontró el género inicial: %v", err)
	}

	cuerpo = `{"titulo":"Pelicula De Prueba Editada","anio":2015,"generoId":` +
		itoa(otroGenero.ID) + `,"estado":"vista","puntuacion":7.5}`
	respuesta = conSesion(t, router, token, http.MethodPut, "/api/peliculas/"+itoa(creada.ID), cuerpo)
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
	}

	var editada models.Pelicula
	decodificar(t, respuesta.Body.Bytes(), &editada)
	if editada.Titulo != "Pelicula De Prueba Editada" || editada.Anio != 2015 {
		t.Fatalf("la película no se modificó correctamente: %+v", editada)
	}
	if editada.GeneroID != otroGenero.ID || editada.Genero.Nombre != "Drama" {
		t.Fatalf("el género no se modificó correctamente: %+v", editada.Genero)
	}
	if editada.Puntuacion == nil || *editada.Puntuacion != 7.5 {
		t.Fatal("la puntuación no se modificó correctamente")
	}

	// No se puede eliminar el género que la película está utilizando.
	respuesta = conSesion(t, router, token, http.MethodDelete, "/api/generos/"+itoa(otroGenero.ID), "")
	if respuesta.Code != http.StatusConflict {
		t.Fatalf("se esperaba 409 para un género en uso, se obtuvo %d", respuesta.Code)
	}

	// Estadísticas.
	respuesta = conSesion(t, router, token, http.MethodGet, "/api/estadisticas", "")
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", respuesta.Code)
	}

	// Eliminar la película.
	respuesta = conSesion(t, router, token, http.MethodDelete, "/api/peliculas/"+itoa(creada.ID), "")
	if respuesta.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d", respuesta.Code)
	}

	// Ya no existe.
	respuesta = conSesion(t, router, token, http.MethodGet, "/api/peliculas/"+itoa(creada.ID), "")
	if respuesta.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", respuesta.Code)
	}
}

func TestIntegracionHealthYGeneros(t *testing.T) {
	router, _, token := prepararIntegracion(t)

	// El healthcheck es público: Docker lo consulta sin token.
	respuesta := ejecutar(t, router, http.MethodGet, "/health", "")
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", respuesta.Code)
	}
	var salud map[string]any
	decodificar(t, respuesta.Body.Bytes(), &salud)
	if salud["status"] != "healthy" {
		t.Fatalf("se esperaba \"healthy\", se obtuvo %v", salud["status"])
	}

	respuesta = conSesion(t, router, token, http.MethodGet, "/api/generos", "")
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", respuesta.Code)
	}
	var generos []models.Genero
	decodificar(t, respuesta.Body.Bytes(), &generos)
	if len(generos) < len(database.GenerosIniciales) {
		t.Fatalf("se esperaban al menos %d géneros, se obtuvieron %d", len(database.GenerosIniciales), len(generos))
	}

	// Los géneros duplicados no están permitidos (ignorando mayúsculas y minúsculas).
	respuesta = conSesion(t, router, token, http.MethodPost, "/api/generos", `{"nombre":"cOMEDIA"}`)
	if respuesta.Code != http.StatusConflict {
		t.Fatalf("se esperaba 409 para un género duplicado, se obtuvo %d", respuesta.Code)
	}
}

// TestIntegracionCuentas comprueba el registro, el login y el perfil contra la base de datos.
func TestIntegracionCuentas(t *testing.T) {
	router, _, token := prepararIntegracion(t)

	// El perfil devuelve la cuenta del token, sin el hash de la contraseña.
	respuesta := conSesion(t, router, token, http.MethodGet, "/api/auth/perfil", "")
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
	}
	var perfil struct {
		Usuario models.Usuario `json:"usuario"`
	}
	decodificar(t, respuesta.Body.Bytes(), &perfil)
	if perfil.Usuario.Email != emailDePrueba {
		t.Fatalf("se esperaba el email %q, se obtuvo %q", emailDePrueba, perfil.Usuario.Email)
	}
	if strings.Contains(respuesta.Body.String(), "PasswordHash") ||
		strings.Contains(respuesta.Body.String(), "$2a$") {
		t.Fatal("la respuesta nunca debe incluir el hash de la contraseña")
	}

	// No se puede registrar dos veces el mismo email, aunque cambie la caja.
	cuerpo := `{"nombre":"Otro","email":"` + strings.ToUpper(emailDePrueba) + `","password":"` + passwordDePrueba + `"}`
	respuesta = ejecutar(t, router, http.MethodPost, "/api/auth/registro", cuerpo)
	if respuesta.Code != http.StatusConflict {
		t.Fatalf("se esperaba 409 para un email repetido, se obtuvo %d", respuesta.Code)
	}

	// Login correcto.
	respuesta = ejecutar(t, router, http.MethodPost, "/api/auth/login",
		`{"email":"`+emailDePrueba+`","password":"`+passwordDePrueba+`"}`)
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
	}
	var sesion struct {
		Token string `json:"token"`
	}
	decodificar(t, respuesta.Body.Bytes(), &sesion)
	if sesion.Token == "" {
		t.Fatal("el login debe devolver un token")
	}

	// El token del login sirve para usar la API.
	respuesta = conSesion(t, router, sesion.Token, http.MethodGet, "/api/peliculas", "")
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200 con el token del login, se obtuvo %d", respuesta.Code)
	}

	// Contraseña incorrecta y email inexistente devuelven el mismo 401.
	respuesta = ejecutar(t, router, http.MethodPost, "/api/auth/login",
		`{"email":"`+emailDePrueba+`","password":"contraseña-incorrecta"}`)
	if respuesta.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401 con la contraseña incorrecta, se obtuvo %d", respuesta.Code)
	}

	respuesta = ejecutar(t, router, http.MethodPost, "/api/auth/login",
		`{"email":"no-existe@ejemplo.com","password":"`+passwordDePrueba+`"}`)
	if respuesta.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401 con un email inexistente, se obtuvo %d", respuesta.Code)
	}
}

// itoa convierte un identificador en texto para armar las rutas.
func itoa(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
