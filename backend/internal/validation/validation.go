// Package validation contiene las reglas de negocio de películas y géneros.
// Son funciones puras (sin base de datos) para que resulten fáciles de testear.
package validation

import (
	"errors"
	"fmt"
	"math"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
)

// Estados permitidos para una película.
const (
	EstadoPendiente = "pendiente"
	EstadoVista     = "vista"
)

// Límites de las validaciones.
const (
	AnioMinimo        = 1888
	AnioMargenFuturo  = 5
	TituloMaxLongitud = 200
	NombreMinLongitud = 2
	NombreMaxLongitud = 50
	PuntuacionMinima  = 1.0
	PuntuacionMaxima  = 10.0

	EmailMaxLongitud    = 120
	PasswordMinLongitud = 8
	PasswordMaxBytes    = 72 // límite de bcrypt
)

// LimpiarTexto elimina los espacios innecesarios al principio y al final.
func LimpiarTexto(texto string) string {
	return strings.TrimSpace(texto)
}

// AnioMaximo devuelve el año máximo permitido (año actual + 5).
// El año actual nunca se hardcodea.
func AnioMaximo() int {
	return time.Now().Year() + AnioMargenFuturo
}

// ValidarTitulo comprueba que el título tenga entre 1 y 200 caracteres.
// Recibe el título ya limpiado con LimpiarTexto.
func ValidarTitulo(titulo string) error {
	if titulo == "" {
		return errors.New("El título es obligatorio.")
	}
	if utf8.RuneCountInString(titulo) > TituloMaxLongitud {
		return fmt.Errorf("El título no puede superar los %d caracteres.", TituloMaxLongitud)
	}
	return nil
}

// ValidarAnio comprueba que el año esté dentro de un rango razonable.
func ValidarAnio(anio int) error {
	maximo := AnioMaximo()
	if anio < AnioMinimo || anio > maximo {
		return fmt.Errorf("El año debe estar entre %d y %d.", AnioMinimo, maximo)
	}
	return nil
}

// ValidarEstado comprueba que el estado sea "pendiente" o "vista".
func ValidarEstado(estado string) error {
	if estado != EstadoPendiente && estado != EstadoVista {
		return errors.New("El estado solamente puede ser \"pendiente\" o \"vista\".")
	}
	return nil
}

// ValidarGeneroID comprueba que se haya indicado un género.
// La existencia del género se comprueba en el handler contra la base de datos.
func ValidarGeneroID(generoID uint) error {
	if generoID == 0 {
		return errors.New("El género es obligatorio.")
	}
	return nil
}

// ValidarPuntuacion comprueba el rango (1 a 10) y que tenga como máximo un decimal.
func ValidarPuntuacion(puntuacion float64) error {
	if math.IsNaN(puntuacion) || math.IsInf(puntuacion, 0) {
		return errors.New("La puntuación no es un número válido.")
	}
	if puntuacion < PuntuacionMinima || puntuacion > PuntuacionMaxima {
		return fmt.Errorf("La puntuación debe estar entre %.0f y %.0f.", PuntuacionMinima, PuntuacionMaxima)
	}
	if math.Abs(puntuacion*10-math.Round(puntuacion*10)) > 1e-9 {
		return errors.New("La puntuación admite como máximo un decimal.")
	}
	return nil
}

// ValidarPuntuacionSegunEstado aplica la regla principal del dominio:
// solamente una película vista puede tener puntuación.
func ValidarPuntuacionSegunEstado(estado string, puntuacion *float64) error {
	if puntuacion == nil {
		return nil
	}
	if estado != EstadoVista {
		return errors.New("Una película pendiente no puede tener puntuación.")
	}
	return ValidarPuntuacion(*puntuacion)
}

// NormalizarPuntuacion devuelve la puntuación que corresponde guardar según el estado.
// Una película pendiente siempre queda con puntuación null.
func NormalizarPuntuacion(estado string, puntuacion *float64) *float64 {
	if estado != EstadoVista {
		return nil
	}
	return puntuacion
}

// ValidarPelicula aplica todas las validaciones de una película y devuelve el primer error.
// El título debe llegar ya limpiado con LimpiarTexto.
func ValidarPelicula(titulo string, anio int, generoID uint, estado string, puntuacion *float64) error {
	if err := ValidarTitulo(titulo); err != nil {
		return err
	}
	if err := ValidarAnio(anio); err != nil {
		return err
	}
	if err := ValidarGeneroID(generoID); err != nil {
		return err
	}
	if err := ValidarEstado(estado); err != nil {
		return err
	}
	return ValidarPuntuacionSegunEstado(estado, puntuacion)
}

// NormalizarEmail limpia el email y lo pasa a minúsculas,
// de manera que no puedan registrarse dos cuentas iguales con distinta caja.
func NormalizarEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ValidarEmail comprueba el formato y la longitud del email.
// Recibe el email ya normalizado con NormalizarEmail.
func ValidarEmail(email string) error {
	if email == "" {
		return errors.New("El email es obligatorio.")
	}
	if len(email) > EmailMaxLongitud {
		return fmt.Errorf("El email no puede superar los %d caracteres.", EmailMaxLongitud)
	}
	direccion, err := mail.ParseAddress(email)
	if err != nil || direccion.Address != email {
		return errors.New("El email no tiene un formato válido.")
	}
	return nil
}

// ValidarPassword comprueba la longitud de la contraseña.
// El máximo son 72 bytes porque es el límite de bcrypt.
func ValidarPassword(password string) error {
	if password == "" {
		return errors.New("La contraseña es obligatoria.")
	}
	if utf8.RuneCountInString(password) < PasswordMinLongitud {
		return fmt.Errorf("La contraseña debe tener al menos %d caracteres.", PasswordMinLongitud)
	}
	if len(password) > PasswordMaxBytes {
		return fmt.Errorf("La contraseña no puede superar los %d bytes.", PasswordMaxBytes)
	}
	return nil
}

// ValidarNombreUsuario comprueba el nombre visible de la cuenta.
// Recibe el nombre ya limpiado con LimpiarTexto.
func ValidarNombreUsuario(nombre string) error {
	longitud := utf8.RuneCountInString(nombre)
	if longitud == 0 {
		return errors.New("El nombre es obligatorio.")
	}
	if longitud < NombreMinLongitud || longitud > NombreMaxLongitud {
		return fmt.Errorf("El nombre debe tener entre %d y %d caracteres.", NombreMinLongitud, NombreMaxLongitud)
	}
	return nil
}

// ValidarNombreGenero comprueba que el nombre tenga entre 2 y 50 caracteres.
// Recibe el nombre ya limpiado con LimpiarTexto.
func ValidarNombreGenero(nombre string) error {
	longitud := utf8.RuneCountInString(nombre)
	if longitud == 0 {
		return errors.New("El nombre del género es obligatorio.")
	}
	if longitud < NombreMinLongitud || longitud > NombreMaxLongitud {
		return fmt.Errorf("El nombre del género debe tener entre %d y %d caracteres.", NombreMinLongitud, NombreMaxLongitud)
	}
	return nil
}
