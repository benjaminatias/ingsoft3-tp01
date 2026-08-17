package tests

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"gestor-peliculas/internal/auth"
	"gestor-peliculas/internal/validation"
)

// TestRutasPrivadasSinToken comprueba que la API no se pueda usar sin iniciar sesión.
func TestRutasPrivadasSinToken(t *testing.T) {
	router := routerSinBaseDeDatos()

	rutas := []struct {
		metodo string
		ruta   string
	}{
		{http.MethodGet, "/api/peliculas"},
		{http.MethodGet, "/api/peliculas/1"},
		{http.MethodPost, "/api/peliculas"},
		{http.MethodPut, "/api/peliculas/1"},
		{http.MethodDelete, "/api/peliculas/1"},
		{http.MethodPatch, "/api/peliculas/1/estado"},
		{http.MethodPatch, "/api/peliculas/1/puntuacion"},
		{http.MethodGet, "/api/generos"},
		{http.MethodPost, "/api/generos"},
		{http.MethodPut, "/api/generos/1"},
		{http.MethodDelete, "/api/generos/1"},
		{http.MethodGet, "/api/estadisticas"},
		{http.MethodGet, "/api/auth/perfil"},
	}

	for _, caso := range rutas {
		t.Run(caso.metodo+" "+caso.ruta, func(t *testing.T) {
			respuesta := ejecutar(t, router, caso.metodo, caso.ruta, "{}")
			if respuesta.Code != http.StatusUnauthorized {
				t.Fatalf("se esperaba 401 sin token, se obtuvo %d", respuesta.Code)
			}

			var cuerpo map[string]any
			if err := json.Unmarshal(respuesta.Body.Bytes(), &cuerpo); err != nil {
				t.Fatalf("la respuesta no es JSON válido: %v", err)
			}
			if mensaje, _ := cuerpo["error"].(string); mensaje == "" {
				t.Fatal("la respuesta 401 debe incluir un mensaje en el campo \"error\"")
			}
		})
	}
}

// TestCabecerasDeAutorizacionInvalidas comprueba el formato del encabezado y la firma.
func TestCabecerasDeAutorizacionInvalidas(t *testing.T) {
	router := routerSinBaseDeDatos()

	// Token firmado con otra clave: no debe aceptarse.
	tokenAjeno := firmarToken(t, []byte("otra-clave-distinta"), jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// Token vencido firmado con la clave correcta.
	tokenVencido := firmarToken(t, []byte(claveDePrueba), jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1",
		"exp": time.Now().Add(-time.Hour).Unix(),
	})

	casos := []struct {
		nombre   string
		cabecera string
	}{
		{"sin encabezado", ""},
		{"esquema incorrecto", "Basic " + tokenDePrueba(t)},
		{"sin el esquema Bearer", tokenDePrueba(t)},
		{"token vacío", "Bearer "},
		{"token con basura", "Bearer esto.no.es.un.token"},
		{"token firmado con otra clave", "Bearer " + tokenAjeno},
		{"token vencido", "Bearer " + tokenVencido},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			respuesta := ejecutarConCabecera(t, router, http.MethodGet, "/api/peliculas", "", caso.cabecera)
			if respuesta.Code != http.StatusUnauthorized {
				t.Fatalf("se esperaba 401, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
			}
		})
	}
}

// TestTokenSinFirmaEsRechazado comprueba que no se acepte el algoritmo "none".
func TestTokenSinFirmaEsRechazado(t *testing.T) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "1",
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("no se pudo generar el token sin firma: %v", err)
	}

	if _, err := auth.ValidarToken(token); err == nil {
		t.Fatal("un token con algoritmo \"none\" nunca debe aceptarse")
	}
}

// TestGenerarYValidarToken comprueba el ciclo completo del token.
func TestGenerarYValidarToken(t *testing.T) {
	token, expira, err := auth.GenerarToken(42, "usuario@ejemplo.com")
	if err != nil {
		t.Fatalf("no se pudo generar el token: %v", err)
	}

	if expira.Before(time.Now()) {
		t.Fatal("el token no debería nacer vencido")
	}

	credenciales, err := auth.ValidarToken(token)
	if err != nil {
		t.Fatalf("el token generado debería ser válido: %v", err)
	}
	if credenciales.UsuarioID != 42 {
		t.Fatalf("se esperaba el usuario 42, se obtuvo %d", credenciales.UsuarioID)
	}
	if credenciales.Email != "usuario@ejemplo.com" {
		t.Fatalf("el email del token no coincide: %s", credenciales.Email)
	}
}

// TestHashearPassword comprueba el hasheo con bcrypt.
func TestHashearPassword(t *testing.T) {
	hash, err := auth.HashearPassword("contraseña-segura")
	if err != nil {
		t.Fatalf("no se pudo hashear la contraseña: %v", err)
	}

	if hash == "contraseña-segura" {
		t.Fatal("la contraseña nunca debe guardarse en texto plano")
	}
	if !auth.VerificarPassword(hash, "contraseña-segura") {
		t.Fatal("la contraseña correcta debería validarse")
	}
	if auth.VerificarPassword(hash, "otra-contraseña") {
		t.Fatal("una contraseña incorrecta no debe validarse")
	}
}

// TestValidacionesDeCuenta comprueba las reglas de registro sin tocar la base de datos.
func TestValidacionesDeCuenta(t *testing.T) {
	if err := validation.ValidarEmail(validation.NormalizarEmail("  Usuario@Ejemplo.COM ")); err != nil {
		t.Fatalf("el email debería ser válido: %v", err)
	}
	if normalizado := validation.NormalizarEmail("  Usuario@Ejemplo.COM "); normalizado != "usuario@ejemplo.com" {
		t.Fatalf("el email debería normalizarse a minúsculas, se obtuvo %q", normalizado)
	}

	emailsInvalidos := []string{"", "sin-arroba", "espacio @ejemplo.com", "dos@@ejemplo.com", "Nombre <n@ejemplo.com>"}
	for _, email := range emailsInvalidos {
		if err := validation.ValidarEmail(validation.NormalizarEmail(email)); err == nil {
			t.Fatalf("se esperaba un error para el email %q", email)
		}
	}

	if err := validation.ValidarPassword("12345678"); err != nil {
		t.Fatalf("una contraseña de 8 caracteres debería ser válida: %v", err)
	}
	if err := validation.ValidarPassword("corta"); err == nil {
		t.Fatal("se esperaba un error para una contraseña de menos de 8 caracteres")
	}
	if err := validation.ValidarPassword(strings.Repeat("a", 73)); err == nil {
		t.Fatal("se esperaba un error para una contraseña de más de 72 bytes (límite de bcrypt)")
	}

	if err := validation.ValidarNombreUsuario("Benja"); err != nil {
		t.Fatalf("el nombre debería ser válido: %v", err)
	}
	if err := validation.ValidarNombreUsuario("B"); err == nil {
		t.Fatal("se esperaba un error para un nombre de un solo carácter")
	}
}

// TestRegistroYLoginInvalidos comprueba los 400 que no necesitan base de datos.
func TestRegistroYLoginInvalidos(t *testing.T) {
	router := routerSinBaseDeDatos()

	casos := []struct {
		nombre string
		ruta   string
		cuerpo string
	}{
		{"registro sin cuerpo", "/api/auth/registro", "{}"},
		{"registro con cuerpo mal formado", "/api/auth/registro", "{no es json}"},
		{"registro con nombre corto", "/api/auth/registro", `{"nombre":"B","email":"b@ejemplo.com","password":"12345678"}`},
		{"registro con email inválido", "/api/auth/registro", `{"nombre":"Benja","email":"sin-arroba","password":"12345678"}`},
		{"registro con contraseña corta", "/api/auth/registro", `{"nombre":"Benja","email":"b@ejemplo.com","password":"corta"}`},
		{"login sin cuerpo", "/api/auth/login", "{}"},
		{"login sin contraseña", "/api/auth/login", `{"email":"b@ejemplo.com","password":""}`},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			respuesta := ejecutar(t, router, http.MethodPost, caso.ruta, caso.cuerpo)
			if respuesta.Code != http.StatusBadRequest {
				t.Fatalf("se esperaba 400, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
			}
		})
	}
}

// firmarToken arma un token a mano para los casos límite.
func firmarToken(t *testing.T, clave []byte, metodo jwt.SigningMethod, datos jwt.MapClaims) string {
	t.Helper()

	token, err := jwt.NewWithClaims(metodo, datos).SignedString(clave)
	if err != nil {
		t.Fatalf("no se pudo firmar el token: %v", err)
	}
	return token
}
