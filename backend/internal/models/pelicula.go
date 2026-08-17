package models

// Pelicula representa una película de la colección personal.
// Estado solamente puede ser "pendiente" o "vista".
// Puntuacion es un puntero porque las películas pendientes no tienen puntuación (null).
type Pelicula struct {
	ID         uint     `json:"id" gorm:"primaryKey"`
	Titulo     string   `json:"titulo" gorm:"size:200;not null;index"`
	Anio       int      `json:"anio" gorm:"not null"`
	GeneroID   uint     `json:"generoId" gorm:"not null;index"`
	Genero     Genero   `json:"genero" gorm:"foreignKey:GeneroID"`
	Estado     string   `json:"estado" gorm:"size:20;not null"`
	Puntuacion *float64 `json:"puntuacion"`
}
