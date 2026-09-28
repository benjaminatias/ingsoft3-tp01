package tests

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gestor-peliculas/internal/models"
)

// Los géneros de prueba llevan un prefijo propio: cada test los borra antes y después,
// así que se pueden correr muchas veces sobre la misma base sin ensuciarla.
const prefijoGeneroPrueba = "ZZ-Test "

func limpiarGenerosDePrueba(db *gorm.DB) {
	db.Where("nombre LIKE ?", prefijoGeneroPrueba+"%").Delete(&models.Genero{})
}

// crearGeneroDePrueba crea un género por la API y devuelve su respuesta ya decodificada.
func crearGeneroDePrueba(t *testing.T, router *gin.Engine, token, sufijo string) models.Genero {
	t.Helper()

	respuesta := conSesion(t, router, token, http.MethodPost, "/api/generos",
		`{"nombre":"`+prefijoGeneroPrueba+sufijo+`"}`)
	if respuesta.Code != http.StatusCreated {
		t.Fatalf("no se pudo crear el género de prueba: %d (%s)", respuesta.Code, respuesta.Body.String())
	}
	var genero models.Genero
	decodificar(t, respuesta.Body.Bytes(), &genero)
	return genero
}

// Ciclo feliz completo de un género: crear, leer, renombrar, listar y borrar.
func TestIntegracionGenerosCicloCompleto(t *testing.T) {
	router, db, token := prepararIntegracion(t)
	limpiarGenerosDePrueba(db)
	t.Cleanup(func() { limpiarGenerosDePrueba(db) })

	genero := crearGeneroDePrueba(t, router, token, "Ciclo")
	ruta := "/api/generos/" + itoa(genero.ID)

	// Leerlo por identificador.
	respuesta := conSesion(t, router, token, http.MethodGet, ruta, "")
	if respuesta.Code != http.StatusOK {
		t.Fatalf("GET: se esperaba 200, se obtuvo %d", respuesta.Code)
	}
	var leido models.Genero
	decodificar(t, respuesta.Body.Bytes(), &leido)
	if leido.Nombre != prefijoGeneroPrueba+"Ciclo" {
		t.Fatalf("GET: nombre inesperado %q", leido.Nombre)
	}

	// Renombrarlo: el espacio de sobra se limpia.
	respuesta = conSesion(t, router, token, http.MethodPut, ruta, `{"nombre":"  `+prefijoGeneroPrueba+`Renombrado  "}`)
	if respuesta.Code != http.StatusOK {
		t.Fatalf("PUT: se esperaba 200, se obtuvo %d (%s)", respuesta.Code, respuesta.Body.String())
	}
	var actualizado models.Genero
	decodificar(t, respuesta.Body.Bytes(), &actualizado)
	if actualizado.Nombre != prefijoGeneroPrueba+"Renombrado" {
		t.Fatalf("PUT: el nombre debía quedar limpio y nuevo, se obtuvo %q", actualizado.Nombre)
	}

	// Aparece en el listado con su nombre nuevo.
	respuesta = conSesion(t, router, token, http.MethodGet, "/api/generos", "")
	var todos []models.Genero
	decodificar(t, respuesta.Body.Bytes(), &todos)
	encontrado := false
	for _, g := range todos {
		if g.ID == genero.ID && g.Nombre == prefijoGeneroPrueba+"Renombrado" {
			encontrado = true
		}
	}
	if !encontrado {
		t.Fatal("el género renombrado debería aparecer en el listado")
	}

	// Borrarlo y comprobar que ya no existe.
	respuesta = conSesion(t, router, token, http.MethodDelete, ruta, "")
	if respuesta.Code != http.StatusNoContent {
		t.Fatalf("DELETE: se esperaba 204, se obtuvo %d", respuesta.Code)
	}
	respuesta = conSesion(t, router, token, http.MethodGet, ruta, "")
	if respuesta.Code != http.StatusNotFound {
		t.Fatalf("GET tras borrar: se esperaba 404, se obtuvo %d", respuesta.Code)
	}
}

// Casos de error de la API de géneros: cada fila es una regla distinta.
func TestIntegracionGenerosErrores(t *testing.T) {
	router, _, token := prepararIntegracion(t)

	casos := []struct {
		nombre   string
		metodo   string
		ruta     string
		cuerpo   string
		esperado int
	}{
		{"crear con cuerpo mal formado", http.MethodPost, "/api/generos", `{no es json}`, http.StatusBadRequest},
		{"crear sin el campo nombre", http.MethodPost, "/api/generos", `{}`, http.StatusBadRequest},
		{"crear con nombre en blanco", http.MethodPost, "/api/generos", `{"nombre":"   "}`, http.StatusBadRequest},
		{"leer con id que no es número", http.MethodGet, "/api/generos/abc", "", http.StatusBadRequest},
		{"leer con id cero", http.MethodGet, "/api/generos/0", "", http.StatusBadRequest},
		{"leer un género inexistente", http.MethodGet, "/api/generos/999999", "", http.StatusNotFound},
		{"renombrar un género inexistente", http.MethodPut, "/api/generos/999999", `{"nombre":"Cualquiera"}`, http.StatusNotFound},
		{"renombrar con nombre inválido", http.MethodPut, "/api/generos/1", `{"nombre":""}`, http.StatusBadRequest},
		{"borrar con id inválido", http.MethodDelete, "/api/generos/abc", "", http.StatusBadRequest},
		{"borrar un género inexistente", http.MethodDelete, "/api/generos/999999", "", http.StatusNotFound},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			respuesta := conSesion(t, router, token, caso.metodo, caso.ruta, caso.cuerpo)
			if respuesta.Code != caso.esperado {
				t.Fatalf("se esperaba %d, se obtuvo %d (%s)", caso.esperado, respuesta.Code, respuesta.Body.String())
			}
		})
	}
}

// Renombrar a un nombre que ya usa OTRO género es un 409, ignorando mayúsculas;
// renombrar a su propio nombre no es un conflicto (por eso se excluye el propio id).
func TestIntegracionGeneroRenombrarDuplicado(t *testing.T) {
	router, db, token := prepararIntegracion(t)
	limpiarGenerosDePrueba(db)
	t.Cleanup(func() { limpiarGenerosDePrueba(db) })

	crearGeneroDePrueba(t, router, token, "Uno")
	segundo := crearGeneroDePrueba(t, router, token, "Dos")
	ruta := "/api/generos/" + itoa(segundo.ID)

	respuesta := conSesion(t, router, token, http.MethodPut, ruta, `{"nombre":"zz-test UNO"}`)
	if respuesta.Code != http.StatusConflict {
		t.Fatalf("se esperaba 409 al usar el nombre de otro género, se obtuvo %d", respuesta.Code)
	}

	respuesta = conSesion(t, router, token, http.MethodPut, ruta, `{"nombre":"`+prefijoGeneroPrueba+`Dos"}`)
	if respuesta.Code != http.StatusOK {
		t.Fatalf("se esperaba 200 al conservar el propio nombre, se obtuvo %d", respuesta.Code)
	}
}

// Un género con películas NO se puede borrar: 409 y el género sigue existiendo.
func TestIntegracionNoSeBorraUnGeneroConPeliculas(t *testing.T) {
	router, db, token := prepararIntegracion(t)
	limpiarGenerosDePrueba(db)

	genero := crearGeneroDePrueba(t, router, token, "Ocupado")

	cuerpo := `{"titulo":"Pelicula De Genero Ocupado","anio":2010,"generoId":` + itoa(genero.ID) + `,"estado":"pendiente"}`
	respuesta := conSesion(t, router, token, http.MethodPost, "/api/peliculas", cuerpo)
	if respuesta.Code != http.StatusCreated {
		t.Fatalf("no se pudo crear la película de prueba: %d (%s)", respuesta.Code, respuesta.Body.String())
	}
	var pelicula models.Pelicula
	decodificar(t, respuesta.Body.Bytes(), &pelicula)

	// Primero se borra la película y después el género: es el orden que exige la clave foránea.
	t.Cleanup(func() {
		db.Unscoped().Delete(&models.Pelicula{}, pelicula.ID)
		limpiarGenerosDePrueba(db)
	})

	ruta := "/api/generos/" + itoa(genero.ID)
	respuesta = conSesion(t, router, token, http.MethodDelete, ruta, "")
	if respuesta.Code != http.StatusConflict {
		t.Fatalf("se esperaba 409 porque el género tiene películas, se obtuvo %d", respuesta.Code)
	}

	respuesta = conSesion(t, router, token, http.MethodGet, ruta, "")
	if respuesta.Code != http.StatusOK {
		t.Fatalf("el género debe seguir existiendo, se obtuvo %d", respuesta.Code)
	}
}
