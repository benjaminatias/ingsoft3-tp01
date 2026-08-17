// Package auth resuelve el hasheo de contraseñas y los tokens JWT.
// No conoce la base de datos ni Gin: son funciones puras, fáciles de testear.
package auth

import (
	"crypto/rand"
	"errors"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// DuracionToken es el tiempo de vida de un token emitido.
const DuracionToken = 24 * time.Hour

// metodoFirma es el único algoritmo aceptado. Fijarlo evita el ataque
// de confusión de algoritmos (por ejemplo un token firmado con "none").
var metodoFirma = jwt.SigningMethodHS256

// ErrTokenInvalido se devuelve para cualquier token que no se pueda utilizar.
var ErrTokenInvalido = errors.New("token inválido o expirado")

var (
	secretoUnaVez sync.Once
	secreto       []byte
)

// claveSecreta lee JWT_SECRET. La clave nunca se hardcodea: si la variable no
// está definida se genera una aleatoria para desarrollo y se avisa por consola
// (al reiniciar el servidor las sesiones anteriores dejan de ser válidas).
func claveSecreta() []byte {
	secretoUnaVez.Do(func() {
		if valor := os.Getenv("JWT_SECRET"); valor != "" {
			secreto = []byte(valor)
			return
		}

		aleatoria := make([]byte, 32)
		if _, err := rand.Read(aleatoria); err != nil {
			log.Fatalf("No se pudo generar una clave para los tokens: %v", err)
		}
		secreto = aleatoria
		log.Println("ADVERTENCIA: JWT_SECRET no está definida. Se generó una clave aleatoria solamente para desarrollo.")
	})
	return secreto
}

// Credenciales son los datos del usuario que viajan dentro del token.
type Credenciales struct {
	UsuarioID uint
	Email     string
}

type datosToken struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerarToken emite un JWT firmado para el usuario indicado.
// Devuelve también la fecha de expiración para informarla al frontend.
func GenerarToken(usuarioID uint, email string) (string, time.Time, error) {
	ahora := time.Now()
	expira := ahora.Add(DuracionToken)

	datos := datosToken{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(usuarioID), 10),
			IssuedAt:  jwt.NewNumericDate(ahora),
			ExpiresAt: jwt.NewNumericDate(expira),
		},
	}

	token, err := jwt.NewWithClaims(metodoFirma, datos).SignedString(claveSecreta())
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expira, nil
}

// ValidarToken comprueba la firma y la expiración del token recibido.
func ValidarToken(cadena string) (Credenciales, error) {
	var datos datosToken

	token, err := jwt.ParseWithClaims(cadena, &datos, func(*jwt.Token) (any, error) {
		return claveSecreta(), nil
	}, jwt.WithValidMethods([]string{metodoFirma.Alg()}))

	if err != nil || !token.Valid {
		return Credenciales{}, ErrTokenInvalido
	}

	usuarioID, err := strconv.ParseUint(datos.Subject, 10, 64)
	if err != nil || usuarioID == 0 {
		return Credenciales{}, ErrTokenInvalido
	}

	return Credenciales{UsuarioID: uint(usuarioID), Email: datos.Email}, nil
}

// HashearPassword genera el hash bcrypt que se guarda en la base de datos.
func HashearPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerificarPassword compara una contraseña con su hash.
func VerificarPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// hashDescartable es un hash válido de una contraseña que nadie utiliza.
// Sirve para gastar el mismo tiempo cuando el email no existe y así no
// revelar, midiendo la demora, qué emails están registrados.
var hashDescartable, _ = bcrypt.GenerateFromPassword([]byte("contraseña-inexistente"), bcrypt.DefaultCost)

// EquilibrarTiempo se llama cuando no se encontró el usuario.
func EquilibrarTiempo() {
	_ = bcrypt.CompareHashAndPassword(hashDescartable, []byte("contraseña-inexistente"))
}
