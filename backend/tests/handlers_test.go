package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"gestor-peliculas/internal/auth"
	"gestor-peliculas/internal/handlers"
)

// claveDePrueba es la clave con la que se firman los tokens durante los tests.
const claveDePrueba = "clave-solamente-para-tests"

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	// Debe definirse antes de generar el primer token: la clave se lee una sola vez.
	os.Setenv("JWT_SECRET", claveDePrueba)
	os.Exit(m.Run())
}

// routerSinBaseDeDatos crea el router real sin conexión.
// Sirve para comprobar las respuestas que no necesitan tocar la base de datos.
func routerSinBaseDeDatos() *gin.Engine {
	return handlers.NuevoRouter(nil)
}

// tokenDePrueba genera un token válido para un usuario cualquiera.
// El middleware solamente valida el token, por eso no hace falta base de datos.
func tokenDePrueba(t *testing.T) string {
	t.Helper()

	token, _, err := auth.GenerarToken(1, "tester@ejemplo.com")
	if err != nil {
		t.Fatalf("no se pudo generar el token de prueba: %v", err)
	}
	return token
}

// ejecutar lanza una petición sin token.
func ejecutar(t *testing.T, router *gin.Engine, metodo, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	return ejecutarConCabecera(t, router, metodo, ruta, cuerpo, "")
}

// ejecutarAutenticado lanza una petición con un token válido.
func ejecutarAutenticado(t *testing.T, router *gin.Engine, metodo, ruta, cuerpo string) *httptest.ResponseRecorder {
	t.Helper()
	return ejecutarConCabecera(t, router, metodo, ruta, cuerpo, "Bearer "+tokenDePrueba(t))
}

func ejecutarConCabecera(t *testing.T, router *gin.Engine, metodo, ruta, cuerpo, autorizacion string) *httptest.ResponseRecorder {
	t.Helper()

	peticion := httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
	peticion.Header.Set("Content-Type", "application/json")
	if autorizacion != "" {
		peticion.Header.Set("Authorization", autorizacion)
	}

	respuesta := httptest.NewRecorder()
	router.ServeHTTP(respuesta, peticion)
	return respuesta
}

// TestHealthSinBaseDeDatos comprueba que /health responda 503 cuando PostgreSQL no está accesible.
func TestHealthSinBaseDeDatos(t *testing.T) {
	respuesta := ejecutar(t, routerSinBaseDeDatos(), http.MethodGet, "/health", "")

	if respuesta.Code != http.StatusServiceUnavailable {
		t.Fatalf("se esperaba 503, se obtuvo %d", respuesta.Code)
	}

	var cuerpo map[string]any
	if err := json.Unmarshal(respuesta.Body.Bytes(), &cuerpo); err != nil {
		t.Fatalf("la respuesta no es JSON válido: %v", err)
	}
	if cuerpo["status"] != "unhealthy" {
		t.Fatalf("se esperaba status \"unhealthy\", se obtuvo %v", cuerpo["status"])
	}
	if _, existe := cuerpo["error"]; !existe {
		t.Fatal("la respuesta de error debe incluir el campo \"error\"")
	}
}

// TestRespuestasInvalidas comprueba las validaciones que responden 400
// antes de llegar a la base de datos. Todas las rutas son privadas,
// por eso se envía un token válido.
func TestRespuestasInvalidas(t *testing.T) {
	router := routerSinBaseDeDatos()

	casos := []struct {
		nombre string
		metodo string
		ruta   string
		cuerpo string
	}{
		{
			nombre: "identificador no numérico",
			metodo: http.MethodGet,
			ruta:   "/api/peliculas/abc",
		},
		{
			nombre: "filtro de estado inválido",
			metodo: http.MethodGet,
			ruta:   "/api/peliculas?estado=viendo",
		},
		{
			nombre: "filtro de año inválido",
			metodo: http.MethodGet,
			ruta:   "/api/peliculas?anio=dosmil",
		},
		{
			nombre: "cuerpo mal formado",
			metodo: http.MethodPost,
			ruta:   "/api/peliculas",
			cuerpo: "{esto no es json}",
		},
		{
			nombre: "faltan campos obligatorios",
			metodo: http.MethodPost,
			ruta:   "/api/peliculas",
			cuerpo: `{"titulo":"Interstellar"}`,
		},
		{
			nombre: "título vacío",
			metodo: http.MethodPost,
			ruta:   "/api/peliculas",
			cuerpo: `{"titulo":"   ","anio":2014,"generoId":1,"estado":"vista","puntuacion":9}`,
		},
		{
			nombre: "año inválido",
			metodo: http.MethodPost,
			ruta:   "/api/peliculas",
			cuerpo: `{"titulo":"Interstellar","anio":1200,"generoId":1,"estado":"vista","puntuacion":9}`,
		},
		{
			nombre: "estado inválido",
			metodo: http.MethodPost,
			ruta:   "/api/peliculas",
			cuerpo: `{"titulo":"Interstellar","anio":2014,"generoId":1,"estado":"viendo"}`,
		},
		{
			nombre: "película pendiente con puntuación",
			metodo: http.MethodPost,
			ruta:   "/api/peliculas",
			cuerpo: `{"titulo":"Dune: Part Two","anio":2024,"generoId":1,"estado":"pendiente","puntuacion":8.5}`,
		},
		{
			nombre: "puntuación mayor a 10",
			metodo: http.MethodPost,
			ruta:   "/api/peliculas",
			cuerpo: `{"titulo":"Interstellar","anio":2014,"generoId":1,"estado":"vista","puntuacion":11}`,
		},
		{
			nombre: "puntuación menor a 1",
			metodo: http.MethodPatch,
			ruta:   "/api/peliculas/1/puntuacion",
			cuerpo: `{"puntuacion":0.5}`,
		},
		{
			nombre: "estado inválido al cambiar estado",
			metodo: http.MethodPatch,
			ruta:   "/api/peliculas/1/estado",
			cuerpo: `{"estado":"terminada"}`,
		},
		{
			nombre: "nombre de género demasiado corto",
			metodo: http.MethodPost,
			ruta:   "/api/generos",
			cuerpo: `{"nombre":"W"}`,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			respuesta := ejecutarAutenticado(t, router, caso.metodo, caso.ruta, caso.cuerpo)

			if respuesta.Code != http.StatusBadRequest {
				t.Fatalf("se esperaba 400, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
			}

			var cuerpo map[string]any
			if err := json.Unmarshal(respuesta.Body.Bytes(), &cuerpo); err != nil {
				t.Fatalf("la respuesta no es JSON válido: %v", err)
			}
			mensaje, existe := cuerpo["error"].(string)
			if !existe || mensaje == "" {
				t.Fatal("la respuesta de error debe incluir un mensaje entendible en el campo \"error\"")
			}
		})
	}
}

// TestRutaInexistente comprueba que una ruta desconocida no devuelva 200.
func TestRutaInexistente(t *testing.T) {
	respuesta := ejecutar(t, routerSinBaseDeDatos(), http.MethodGet, "/api/inexistente", "")
	if respuesta.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", respuesta.Code)
	}
}
