// Package database se encarga de la conexión con PostgreSQL, las migraciones
// y la creación de los géneros iniciales.
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gestor-peliculas/internal/models"
)

// GenerosIniciales son los géneros que se crean automáticamente la primera vez.
var GenerosIniciales = []string{
	"Acción",
	"Animación",
	"Aventura",
	"Ciencia ficción",
	"Comedia",
	"Crimen",
	"Documental",
	"Drama",
	"Fantasía",
	"Terror",
	"Thriller",
	"Otros",
}

// variableEntorno devuelve el valor de una variable de entorno o un valor por defecto.
// Las credenciales nunca se hardcodean: la contraseña no tiene valor por defecto.
func variableEntorno(clave, porDefecto string) string {
	if valor := os.Getenv(clave); valor != "" {
		return valor
	}
	return porDefecto
}

// cadenaConexion arma el DSN de PostgreSQL a partir de las variables de entorno.
func cadenaConexion() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		variableEntorno("DB_HOST", "localhost"),
		variableEntorno("DB_PORT", "5432"),
		variableEntorno("DB_USER", "postgres"),
		os.Getenv("DB_PASSWORD"),
		variableEntorno("DB_NAME", "peliculas_db"),
		variableEntorno("DB_SSLMODE", "disable"),
	)
}

// Conectar abre la conexión con PostgreSQL reintentando algunas veces,
// porque en Docker la base puede tardar un momento en aceptar conexiones.
func Conectar() (*gorm.DB, error) {
	dsn := cadenaConexion()
	configuracion := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}

	var ultimoError error
	for intento := 1; intento <= 10; intento++ {
		db, err := gorm.Open(postgres.Open(dsn), configuracion)
		if err == nil {
			sqlDB, errSQL := db.DB()
			if errSQL == nil {
				if errPing := sqlDB.Ping(); errPing == nil {
					return db, nil
				} else {
					ultimoError = errPing
				}
			} else {
				ultimoError = errSQL
			}
		} else {
			ultimoError = err
		}

		log.Printf("No se pudo conectar con PostgreSQL (intento %d/10), reintentando...", intento)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("no se pudo conectar con PostgreSQL: %w", ultimoError)
}

// Migrar crea o actualiza las tablas de la base de datos.
func Migrar(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Genero{},
		&models.Pelicula{},
		&models.Usuario{},
	)
}

// CrearGenerosIniciales inserta los géneros base si todavía no existen.
func CrearGenerosIniciales(db *gorm.DB) error {
	for _, nombre := range GenerosIniciales {
		var cantidad int64
		if err := db.Model(&models.Genero{}).
			Where("LOWER(nombre) = LOWER(?)", nombre).
			Count(&cantidad).Error; err != nil {
			return err
		}
		if cantidad > 0 {
			continue
		}
		if err := db.Create(&models.Genero{Nombre: nombre}).Error; err != nil {
			return err
		}
	}
	return nil
}
