package handlers

import (
	"errors"

	"gorm.io/gorm"

	"gestor-peliculas/internal/models"
)

// ErrUsuarioNoEncontrado lo devuelve cualquier repositorio cuando la cuenta no existe.
// Así los handlers no dependen del error propio de GORM.
var ErrUsuarioNoEncontrado = errors.New("usuario no encontrado")

// RepositorioUsuarios es lo único que los handlers de cuentas necesitan del almacenamiento.
// Es una interfaz para que la dependencia entre desde afuera: en producción la implementa
// GORM y en los tests se reemplaza por un doble, sin base de datos.
type RepositorioUsuarios interface {
	ExisteEmail(email string) (bool, error)
	Crear(usuario *models.Usuario) error
	BuscarPorEmail(email string) (models.Usuario, error)
	BuscarPorID(id uint) (models.Usuario, error)
}

// usuariosGorm es la implementación real, sobre PostgreSQL.
type usuariosGorm struct {
	db *gorm.DB
}

func (r usuariosGorm) ExisteEmail(email string) (bool, error) {
	var cantidad int64
	err := r.db.Model(&models.Usuario{}).Where("email = ?", email).Count(&cantidad).Error
	return cantidad > 0, err
}

func (r usuariosGorm) Crear(usuario *models.Usuario) error {
	return r.db.Create(usuario).Error
}

func (r usuariosGorm) BuscarPorEmail(email string) (models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.Where("email = ?", email).First(&usuario).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return usuario, ErrUsuarioNoEncontrado
	}
	return usuario, err
}

func (r usuariosGorm) BuscarPorID(id uint) (models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.First(&usuario, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return usuario, ErrUsuarioNoEncontrado
	}
	return usuario, err
}
