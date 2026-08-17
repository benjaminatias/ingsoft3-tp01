package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gestor-peliculas/internal/auth"
	"gestor-peliculas/internal/models"
	"gestor-peliculas/internal/validation"
)

// claveUsuario es donde el middleware deja el usuario autenticado del contexto.
const claveUsuario = "usuarioID"

type peticionRegistro struct {
	Nombre   *string `json:"nombre"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type peticionLogin struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

// Registro crea una cuenta nueva y devuelve un token para iniciar sesión directamente.
func (h *Handler) Registro(c *gin.Context) {
	var peticion peticionRegistro
	if err := c.ShouldBindJSON(&peticion); err != nil {
		responderError(c, http.StatusBadRequest, "El cuerpo de la petición no es válido.")
		return
	}
	if peticion.Nombre == nil || peticion.Email == nil || peticion.Password == nil {
		responderError(c, http.StatusBadRequest, "Debe indicar nombre, email y contraseña.")
		return
	}

	nombre := validation.LimpiarTexto(*peticion.Nombre)
	email := validation.NormalizarEmail(*peticion.Email)

	if err := validation.ValidarNombreUsuario(nombre); err != nil {
		responderError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validation.ValidarEmail(email); err != nil {
		responderError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validation.ValidarPassword(*peticion.Password); err != nil {
		responderError(c, http.StatusBadRequest, err.Error())
		return
	}

	var cantidad int64
	if err := h.DB.Model(&models.Usuario{}).Where("email = ?", email).Count(&cantidad).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo crear la cuenta.")
		return
	}
	if cantidad > 0 {
		responderError(c, http.StatusConflict, "Ya existe una cuenta con ese email.")
		return
	}

	hash, err := auth.HashearPassword(*peticion.Password)
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo crear la cuenta.")
		return
	}

	usuario := models.Usuario{Nombre: nombre, Email: email, PasswordHash: hash}
	if err := h.DB.Create(&usuario).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo crear la cuenta.")
		return
	}

	responderConToken(c, http.StatusCreated, usuario)
}

// Login valida las credenciales y devuelve un token.
func (h *Handler) Login(c *gin.Context) {
	var peticion peticionLogin
	if err := c.ShouldBindJSON(&peticion); err != nil {
		responderError(c, http.StatusBadRequest, "El cuerpo de la petición no es válido.")
		return
	}
	if peticion.Email == nil || peticion.Password == nil {
		responderError(c, http.StatusBadRequest, "Debe indicar email y contraseña.")
		return
	}

	email := validation.NormalizarEmail(*peticion.Email)
	if email == "" || *peticion.Password == "" {
		responderError(c, http.StatusBadRequest, "Debe indicar email y contraseña.")
		return
	}

	var usuario models.Usuario
	err := h.DB.Where("email = ?", email).First(&usuario).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Se compara igualmente contra un hash descartable para que el tiempo de
		// respuesta no revele si el email existe.
		auth.EquilibrarTiempo()
		responderError(c, http.StatusUnauthorized, "Email o contraseña incorrectos.")
		return
	}
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo iniciar sesión.")
		return
	}

	if !auth.VerificarPassword(usuario.PasswordHash, *peticion.Password) {
		responderError(c, http.StatusUnauthorized, "Email o contraseña incorrectos.")
		return
	}

	responderConToken(c, http.StatusOK, usuario)
}

// Perfil devuelve los datos del usuario autenticado.
// El frontend lo utiliza al abrir la aplicación para saber si el token sigue siendo válido.
func (h *Handler) Perfil(c *gin.Context) {
	usuarioID, ok := UsuarioAutenticado(c)
	if !ok {
		responderError(c, http.StatusUnauthorized, "Debe iniciar sesión.")
		return
	}

	var usuario models.Usuario
	err := h.DB.First(&usuario, usuarioID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responderError(c, http.StatusUnauthorized, "La cuenta ya no existe.")
		return
	}
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo obtener el perfil.")
		return
	}

	c.JSON(http.StatusOK, gin.H{"usuario": usuario})
}

// responderConToken emite el token y devuelve los datos públicos del usuario.
func responderConToken(c *gin.Context, estado int, usuario models.Usuario) {
	token, expira, err := auth.GenerarToken(usuario.ID, usuario.Email)
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo generar el token de sesión.")
		return
	}

	c.JSON(estado, gin.H{
		"token":   token,
		"expira":  expira,
		"usuario": usuario,
	})
}

// RequiereAutenticacion es el middleware que protege las rutas privadas.
// Solamente comprueba el token: no consulta la base de datos.
func RequiereAutenticacion() gin.HandlerFunc {
	return func(c *gin.Context) {
		cabecera := c.GetHeader("Authorization")
		partes := strings.Fields(cabecera)
		if len(partes) != 2 || !strings.EqualFold(partes[0], "Bearer") {
			responderError(c, http.StatusUnauthorized, "Debe iniciar sesión.")
			c.Abort()
			return
		}

		credenciales, err := auth.ValidarToken(partes[1])
		if err != nil {
			responderError(c, http.StatusUnauthorized, "La sesión expiró o el token no es válido.")
			c.Abort()
			return
		}

		c.Set(claveUsuario, credenciales.UsuarioID)
		c.Next()
	}
}

// UsuarioAutenticado devuelve el identificador del usuario de la petición actual.
func UsuarioAutenticado(c *gin.Context) (uint, bool) {
	valor, existe := c.Get(claveUsuario)
	if !existe {
		return 0, false
	}
	usuarioID, ok := valor.(uint)
	return usuarioID, ok
}
