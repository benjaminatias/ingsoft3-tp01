package tests

import (
	"testing"
	"time"

	"gestor-peliculas/internal/validation"
)

func puntero(valor float64) *float64 {
	return &valor
}

func TestTituloVacioEsInvalido(t *testing.T) {
	if err := validation.ValidarTitulo(""); err == nil {
		t.Fatal("se esperaba un error para un título vacío")
	}

	// Un título con solamente espacios queda vacío tras limpiarlo.
	if err := validation.ValidarTitulo(validation.LimpiarTexto("   ")); err == nil {
		t.Fatal("se esperaba un error para un título con solamente espacios")
	}
}

func TestTituloValidoYLimpieza(t *testing.T) {
	titulo := validation.LimpiarTexto("  Interstellar  ")
	if titulo != "Interstellar" {
		t.Fatalf("se esperaba \"Interstellar\", se obtuvo %q", titulo)
	}
	if err := validation.ValidarTitulo(titulo); err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}
}

func TestTituloDemasiadoLargoEsInvalido(t *testing.T) {
	largo := make([]rune, validation.TituloMaxLongitud+1)
	for i := range largo {
		largo[i] = 'a'
	}
	if err := validation.ValidarTitulo(string(largo)); err == nil {
		t.Fatal("se esperaba un error para un título de más de 200 caracteres")
	}
}

func TestAnioInvalido(t *testing.T) {
	casos := []int{0, 1800, 1887, time.Now().Year() + 6}
	for _, anio := range casos {
		if err := validation.ValidarAnio(anio); err == nil {
			t.Fatalf("se esperaba un error para el año %d", anio)
		}
	}
}

func TestAnioValido(t *testing.T) {
	casos := []int{1888, 2014, 2024, time.Now().Year(), time.Now().Year() + 5}
	for _, anio := range casos {
		if err := validation.ValidarAnio(anio); err != nil {
			t.Fatalf("no se esperaba un error para el año %d: %v", anio, err)
		}
	}
}

func TestEstadoInvalido(t *testing.T) {
	casos := []string{"", "VISTA", "viendo", "terminada", "pendientes"}
	for _, estado := range casos {
		if err := validation.ValidarEstado(estado); err == nil {
			t.Fatalf("se esperaba un error para el estado %q", estado)
		}
	}
}

func TestEstadosValidos(t *testing.T) {
	for _, estado := range []string{validation.EstadoPendiente, validation.EstadoVista} {
		if err := validation.ValidarEstado(estado); err != nil {
			t.Fatalf("no se esperaba un error para el estado %q: %v", estado, err)
		}
	}
}

func TestPuntuacionMenorAUnoEsInvalida(t *testing.T) {
	for _, valor := range []float64{0, 0.5, -3} {
		if err := validation.ValidarPuntuacion(valor); err == nil {
			t.Fatalf("se esperaba un error para la puntuación %v", valor)
		}
	}
}

func TestPuntuacionMayorADiezEsInvalida(t *testing.T) {
	for _, valor := range []float64{10.1, 11, 100} {
		if err := validation.ValidarPuntuacion(valor); err == nil {
			t.Fatalf("se esperaba un error para la puntuación %v", valor)
		}
	}
}

func TestPuntuacionConMasDeUnDecimalEsInvalida(t *testing.T) {
	if err := validation.ValidarPuntuacion(8.55); err == nil {
		t.Fatal("se esperaba un error para una puntuación con dos decimales")
	}
}

func TestPuntuacionesValidas(t *testing.T) {
	for _, valor := range []float64{1, 6, 7.5, 8, 8.5, 9.5, 10} {
		if err := validation.ValidarPuntuacion(valor); err != nil {
			t.Fatalf("no se esperaba un error para la puntuación %v: %v", valor, err)
		}
	}
}

func TestPeliculaPendienteConPuntuacionEsInvalida(t *testing.T) {
	err := validation.ValidarPelicula("Dune: Part Two", 2024, 4, validation.EstadoPendiente, puntero(8.5))
	if err == nil {
		t.Fatal("se esperaba un error: una película pendiente no puede tener puntuación")
	}
}

func TestPeliculaPendienteSinPuntuacionEsValida(t *testing.T) {
	if err := validation.ValidarPelicula("Dune: Part Two", 2024, 4, validation.EstadoPendiente, nil); err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}
}

func TestPeliculaVistaConPuntuacionValida(t *testing.T) {
	if err := validation.ValidarPelicula("Interstellar", 2014, 4, validation.EstadoVista, puntero(9.5)); err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}
}

func TestPeliculaSinGeneroEsInvalida(t *testing.T) {
	if err := validation.ValidarPelicula("Interstellar", 2014, 0, validation.EstadoVista, puntero(9.5)); err == nil {
		t.Fatal("se esperaba un error: el género es obligatorio")
	}
}

func TestPasarAPendienteEliminaLaPuntuacion(t *testing.T) {
	puntuacion := puntero(9.5)

	if resultado := validation.NormalizarPuntuacion(validation.EstadoVista, puntuacion); resultado == nil || *resultado != 9.5 {
		t.Fatal("una película vista debe conservar su puntuación")
	}

	if resultado := validation.NormalizarPuntuacion(validation.EstadoPendiente, puntuacion); resultado != nil {
		t.Fatalf("al pasar a pendiente la puntuación debe quedar en null, se obtuvo %v", *resultado)
	}
}

func TestValidarNombreGenero(t *testing.T) {
	invalidos := []string{"", "A", string(make([]rune, validation.NombreMaxLongitud+1))}
	for _, nombre := range invalidos {
		if err := validation.ValidarNombreGenero(nombre); err == nil {
			t.Fatalf("se esperaba un error para el nombre %q", nombre)
		}
	}

	if err := validation.ValidarNombreGenero("Western"); err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}
}
