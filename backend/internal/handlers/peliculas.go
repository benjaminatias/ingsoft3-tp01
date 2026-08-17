package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gestor-peliculas/internal/models"
	"gestor-peliculas/internal/validation"
)

// peticionPelicula representa el cuerpo recibido al crear o modificar una película.
// Se utilizan punteros para distinguir un campo ausente de un valor vacío.
type peticionPelicula struct {
	Titulo     *string  `json:"titulo"`
	Anio       *int     `json:"anio"`
	GeneroID   *uint    `json:"generoId"`
	Estado     *string  `json:"estado"`
	Puntuacion *float64 `json:"puntuacion"`
}

type peticionEstado struct {
	Estado *string `json:"estado"`
}

type peticionPuntuacion struct {
	Puntuacion *float64 `json:"puntuacion"`
}

// filtrosPeliculas contiene los filtros opcionales del listado ya validados.
type filtrosPeliculas struct {
	Estado        string
	GeneroID      uint
	Anio          *int
	Titulo        string
	PuntuacionMin *float64
}

// leerFiltros lee y valida los parámetros de la query string.
// Si alguno no es válido responde 400 y devuelve ok = false.
func leerFiltros(c *gin.Context) (filtrosPeliculas, bool) {
	var filtros filtrosPeliculas

	if estado := c.Query("estado"); estado != "" {
		if err := validation.ValidarEstado(estado); err != nil {
			responderError(c, http.StatusBadRequest, err.Error())
			return filtros, false
		}
		filtros.Estado = estado
	}

	if valor := c.Query("generoId"); valor != "" {
		generoID, err := strconv.ParseUint(valor, 10, 64)
		if err != nil || generoID == 0 {
			responderError(c, http.StatusBadRequest, "El filtro de género no es válido.")
			return filtros, false
		}
		filtros.GeneroID = uint(generoID)
	}

	if valor := c.Query("anio"); valor != "" {
		anio, err := strconv.Atoi(valor)
		if err != nil {
			responderError(c, http.StatusBadRequest, "El filtro de año no es válido.")
			return filtros, false
		}
		filtros.Anio = &anio
	}

	filtros.Titulo = validation.LimpiarTexto(c.Query("titulo"))

	if valor := c.Query("puntuacionMin"); valor != "" {
		minima, err := strconv.ParseFloat(valor, 64)
		if err != nil {
			responderError(c, http.StatusBadRequest, "El filtro de puntuación mínima no es válido.")
			return filtros, false
		}
		filtros.PuntuacionMin = &minima
	}

	return filtros, true
}

// ListarPeliculas devuelve la colección, con el género incluido y ordenada por título.
// Admite filtros opcionales combinables: estado, generoId, anio, titulo y puntuacionMin.
func (h *Handler) ListarPeliculas(c *gin.Context) {
	filtros, ok := leerFiltros(c)
	if !ok {
		return
	}

	consulta := h.DB.Preload("Genero")
	if filtros.Estado != "" {
		consulta = consulta.Where("estado = ?", filtros.Estado)
	}
	if filtros.GeneroID != 0 {
		consulta = consulta.Where("genero_id = ?", filtros.GeneroID)
	}
	if filtros.Anio != nil {
		consulta = consulta.Where("anio = ?", *filtros.Anio)
	}
	if filtros.Titulo != "" {
		consulta = consulta.Where("titulo ILIKE ?", "%"+filtros.Titulo+"%")
	}
	if filtros.PuntuacionMin != nil {
		consulta = consulta.Where("puntuacion >= ?", *filtros.PuntuacionMin)
	}

	peliculas := []models.Pelicula{}
	if err := consulta.Order("LOWER(titulo) ASC").Find(&peliculas).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudieron obtener las películas.")
		return
	}

	c.JSON(http.StatusOK, peliculas)
}

// ObtenerPelicula devuelve una película por su identificador.
func (h *Handler) ObtenerPelicula(c *gin.Context) {
	id, ok := leerID(c)
	if !ok {
		return
	}

	pelicula, encontrada, err := h.buscarPelicula(id)
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo obtener la película.")
		return
	}
	if !encontrada {
		responderError(c, http.StatusNotFound, "La película no existe.")
		return
	}

	c.JSON(http.StatusOK, pelicula)
}

// CrearPelicula registra una película nueva.
func (h *Handler) CrearPelicula(c *gin.Context) {
	var peticion peticionPelicula
	if err := c.ShouldBindJSON(&peticion); err != nil {
		responderError(c, http.StatusBadRequest, "El cuerpo de la petición no es válido.")
		return
	}

	if peticion.Titulo == nil || peticion.Anio == nil || peticion.GeneroID == nil || peticion.Estado == nil {
		responderError(c, http.StatusBadRequest, "Debe indicar título, año, género y estado.")
		return
	}

	titulo := validation.LimpiarTexto(*peticion.Titulo)
	estado := validation.LimpiarTexto(*peticion.Estado)

	if err := validation.ValidarPelicula(titulo, *peticion.Anio, *peticion.GeneroID, estado, peticion.Puntuacion); err != nil {
		responderError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !h.generoExiste(c, *peticion.GeneroID) {
		return
	}

	pelicula := models.Pelicula{
		Titulo:     titulo,
		Anio:       *peticion.Anio,
		GeneroID:   *peticion.GeneroID,
		Estado:     estado,
		Puntuacion: validation.NormalizarPuntuacion(estado, peticion.Puntuacion),
	}

	if err := h.DB.Omit("Genero").Create(&pelicula).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo crear la película.")
		return
	}

	h.responderPelicula(c, http.StatusCreated, pelicula.ID)
}

// ActualizarPelicula modifica título, año, género, estado y puntuación.
func (h *Handler) ActualizarPelicula(c *gin.Context) {
	id, ok := leerID(c)
	if !ok {
		return
	}

	var peticion peticionPelicula
	if err := c.ShouldBindJSON(&peticion); err != nil {
		responderError(c, http.StatusBadRequest, "El cuerpo de la petición no es válido.")
		return
	}

	if peticion.Titulo == nil || peticion.Anio == nil || peticion.GeneroID == nil || peticion.Estado == nil {
		responderError(c, http.StatusBadRequest, "Debe indicar título, año, género y estado.")
		return
	}

	titulo := validation.LimpiarTexto(*peticion.Titulo)
	estado := validation.LimpiarTexto(*peticion.Estado)

	if err := validation.ValidarPelicula(titulo, *peticion.Anio, *peticion.GeneroID, estado, peticion.Puntuacion); err != nil {
		responderError(c, http.StatusBadRequest, err.Error())
		return
	}

	pelicula, encontrada, err := h.buscarPelicula(id)
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo actualizar la película.")
		return
	}
	if !encontrada {
		responderError(c, http.StatusNotFound, "La película no existe.")
		return
	}

	if !h.generoExiste(c, *peticion.GeneroID) {
		return
	}

	pelicula.Titulo = titulo
	pelicula.Anio = *peticion.Anio
	pelicula.GeneroID = *peticion.GeneroID
	pelicula.Estado = estado
	pelicula.Puntuacion = validation.NormalizarPuntuacion(estado, peticion.Puntuacion)

	if err := h.guardarPelicula(&pelicula); err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo actualizar la película.")
		return
	}

	h.responderPelicula(c, http.StatusOK, pelicula.ID)
}

// EliminarPelicula borra una película de la colección.
func (h *Handler) EliminarPelicula(c *gin.Context) {
	id, ok := leerID(c)
	if !ok {
		return
	}

	resultado := h.DB.Delete(&models.Pelicula{}, id)
	if resultado.Error != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo eliminar la película.")
		return
	}
	if resultado.RowsAffected == 0 {
		responderError(c, http.StatusNotFound, "La película no existe.")
		return
	}

	c.Status(http.StatusNoContent)
}

// CambiarEstado marca una película como vista o pendiente.
// Al pasar a pendiente la puntuación se elimina automáticamente.
func (h *Handler) CambiarEstado(c *gin.Context) {
	id, ok := leerID(c)
	if !ok {
		return
	}

	var peticion peticionEstado
	if err := c.ShouldBindJSON(&peticion); err != nil {
		responderError(c, http.StatusBadRequest, "El cuerpo de la petición no es válido.")
		return
	}
	if peticion.Estado == nil {
		responderError(c, http.StatusBadRequest, "Debe indicar el estado.")
		return
	}

	estado := validation.LimpiarTexto(*peticion.Estado)
	if err := validation.ValidarEstado(estado); err != nil {
		responderError(c, http.StatusBadRequest, err.Error())
		return
	}

	pelicula, encontrada, err := h.buscarPelicula(id)
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo cambiar el estado de la película.")
		return
	}
	if !encontrada {
		responderError(c, http.StatusNotFound, "La película no existe.")
		return
	}

	pelicula.Estado = estado
	pelicula.Puntuacion = validation.NormalizarPuntuacion(estado, pelicula.Puntuacion)

	if err := h.guardarPelicula(&pelicula); err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo cambiar el estado de la película.")
		return
	}

	h.responderPelicula(c, http.StatusOK, pelicula.ID)
}

// PuntuarPelicula asigna una puntuación, solamente si la película está vista.
func (h *Handler) PuntuarPelicula(c *gin.Context) {
	id, ok := leerID(c)
	if !ok {
		return
	}

	var peticion peticionPuntuacion
	if err := c.ShouldBindJSON(&peticion); err != nil {
		responderError(c, http.StatusBadRequest, "El cuerpo de la petición no es válido.")
		return
	}
	if peticion.Puntuacion == nil {
		responderError(c, http.StatusBadRequest, "Debe indicar la puntuación.")
		return
	}
	if err := validation.ValidarPuntuacion(*peticion.Puntuacion); err != nil {
		responderError(c, http.StatusBadRequest, err.Error())
		return
	}

	pelicula, encontrada, err := h.buscarPelicula(id)
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo puntuar la película.")
		return
	}
	if !encontrada {
		responderError(c, http.StatusNotFound, "La película no existe.")
		return
	}

	if pelicula.Estado != validation.EstadoVista {
		responderError(c, http.StatusBadRequest, "Solamente se puede puntuar una película marcada como vista.")
		return
	}

	pelicula.Puntuacion = peticion.Puntuacion
	if err := h.guardarPelicula(&pelicula); err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo puntuar la película.")
		return
	}

	h.responderPelicula(c, http.StatusOK, pelicula.ID)
}

// buscarPelicula devuelve una película con su género.
// El segundo valor indica si fue encontrada.
func (h *Handler) buscarPelicula(id uint) (models.Pelicula, bool, error) {
	var pelicula models.Pelicula
	err := h.DB.Preload("Genero").First(&pelicula, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return pelicula, false, nil
	}
	if err != nil {
		return pelicula, false, err
	}
	return pelicula, true, nil
}

// guardarPelicula persiste los cambios, incluyendo la puntuación cuando queda en null.
// Se omite la asociación Genero: como viene precargada, GORM volvería a escribir
// el genero_id anterior y el cambio de género se perdería.
func (h *Handler) guardarPelicula(pelicula *models.Pelicula) error {
	return h.DB.Model(pelicula).Omit("Genero").Updates(map[string]any{
		"titulo":     pelicula.Titulo,
		"anio":       pelicula.Anio,
		"genero_id":  pelicula.GeneroID,
		"estado":     pelicula.Estado,
		"puntuacion": pelicula.Puntuacion,
	}).Error
}

// generoExiste comprueba que el género indicado exista y responde 400 si no.
func (h *Handler) generoExiste(c *gin.Context, generoID uint) bool {
	var cantidad int64
	if err := h.DB.Model(&models.Genero{}).Where("id = ?", generoID).Count(&cantidad).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo comprobar el género indicado.")
		return false
	}
	if cantidad == 0 {
		responderError(c, http.StatusBadRequest, "El género indicado no existe.")
		return false
	}
	return true
}

// responderPelicula devuelve la película con su género ya cargado.
func (h *Handler) responderPelicula(c *gin.Context, estado int, id uint) {
	pelicula, encontrada, err := h.buscarPelicula(id)
	if err != nil || !encontrada {
		responderError(c, http.StatusInternalServerError, "No se pudo obtener la película.")
		return
	}
	c.JSON(estado, pelicula)
}
