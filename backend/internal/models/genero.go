package models

// Genero representa una categoría de películas.
// Un género puede tener muchas películas (relación 1:N).
type Genero struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre" gorm:"size:50;not null;uniqueIndex"`
}
