package models

import "time"

// Usuario representa una cuenta que puede iniciar sesión en la aplicación.
// PasswordHash nunca se serializa a JSON: el hash no debe salir del backend.
type Usuario struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Nombre       string    `json:"nombre" gorm:"size:50;not null"`
	Email        string    `json:"email" gorm:"size:120;not null;uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"size:100;not null"`
	CreadoEn     time.Time `json:"creadoEn" gorm:"autoCreateTime"`
}
