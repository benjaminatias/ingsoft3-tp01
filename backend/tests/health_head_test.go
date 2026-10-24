package tests

import (
	"net/http"
	"testing"
)

// El monitor externo (plan gratuito de UptimeRobot) pregunta con HEAD. Si el endpoint solo contestara
// GET, el monitor vería un 404 y marcaría la app caída aunque estuviera viva.
// Sin base de datos el /health responde 503: lo que importa acá es que la RUTA existe para HEAD
// (un 404 o un 405 querría decir que HEAD no está registrado).
func TestHealthRespondeAHEADIgualQueAGET(t *testing.T) {
	router := routerSinBaseDeDatos()

	for _, metodo := range []string{http.MethodGet, http.MethodHead} {
		respuesta := ejecutar(t, router, metodo, "/health", "")

		if respuesta.Code == http.StatusNotFound || respuesta.Code == http.StatusMethodNotAllowed {
			t.Fatalf("%s /health: la ruta no está registrada para ese método (código %d)", metodo, respuesta.Code)
		}
		if respuesta.Code != http.StatusServiceUnavailable {
			t.Fatalf("%s /health sin base: se esperaba 503, se obtuvo %d", metodo, respuesta.Code)
		}
	}
}
