package models

import "time"

// Definição da estrutura para os espaços reserváveis
type Space struct {
	ID               int       `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	MaxReservations  int       `json:"max_reservations" db:"max_reservations"`
	MaxOccupancy     int       `json:"max_occupancy" db:"max_occupancy"`
	Price            float64   `json:"price" db:"price"`
	AvailableFrom    string    `json:"available_from" db:"available_from"`   // Formato HH:MM
	AvailableUntil   string    `json:"available_until" db:"available_until"` // Formato HH:MM
	Amenities        string    `json:"amenities" db:"amenities"`
	RequiresApproval bool      `json:"requires_approval" db:"requires_approval"`
	ReservationRules string    `json:"reservation_rules" db:"reservation_rules"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	ImageURL         string    `json:"image_url" db:"image_url"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
