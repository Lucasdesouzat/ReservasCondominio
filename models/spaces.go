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

// Função para retornar configurações padrão para cada espaço
func GetDefaultSpaces() []Space {
	return []Space{
		{
			ID:               1,
			Name:             "Salão de Festas",
			Description:      "Espaço para eventos e festas, capacidade de 1 reserva por dia.",
			MaxReservations:  1,
			MaxOccupancy:     1,
			Price:            100.00,
			AvailableFrom:    "10:00",
			AvailableUntil:   "00:00",
			ReservationRules: "Horário fixo das 10h às 00h. Não é necessário informar o horário de início.",
			IsActive:         true,
		},
		{
			ID:               2,
			Name:             "Brinquedoteca",
			Description:      "Espaço infantil com brinquedos variados, capacidade de 1 reserva por vez.",
			MaxReservations:  1,
			MaxOccupancy:     1,
			Price:            0.00,
			AvailableFrom:    "08:00",
			AvailableUntil:   "22:00",
			ReservationRules: "Reservas apenas em horários inteiros, com duração de 1 hora.",
			IsActive:         true,
		},
		{
			ID:               3,
			Name:             "Academia",
			Description:      "Espaço para atividades físicas e exercícios, disponível 24 horas.",
			MaxReservations:  1,
			MaxOccupancy:     4,
			Price:            0.00,
			AvailableFrom:    "00:00",
			AvailableUntil:   "23:59",
			ReservationRules: "Reservas de 1 hora, disponíveis 24 horas.",
			IsActive:         true,
		},
		{
			ID:               4,
			Name:             "Quiosque 1",
			Description:      "Quiosque para uso recreativo, com capacidade para 1 reserva por vez.",
			MaxReservations:  1,
			MaxOccupancy:     1,
			Price:            10.00,
			AvailableFrom:    "10:00",
			AvailableUntil:   "22:00",
			ReservationRules: "Reservas para manhã (10h-16h) ou tarde (16h-22h).",
			IsActive:         true,
		},
		{
			ID:               5,
			Name:             "Quiosque 2",
			Description:      "Quiosque para uso recreativo, com capacidade para 1 reserva por vez.",
			MaxReservations:  1,
			MaxOccupancy:     1,
			Price:            10.00,
			AvailableFrom:    "10:00",
			AvailableUntil:   "22:00",
			ReservationRules: "Reservas para manhã (10h-16h) ou tarde (16h-22h).",
			IsActive:         true,
		},
		{
			ID:               6,
			Name:             "Piscina",
			Description:      "Piscina do condomínio, capacidade para 8 pessoas por vez.",
			MaxReservations:  8,
			MaxOccupancy:     8,
			Price:            0.00,
			AvailableFrom:    "09:00",
			AvailableUntil:   "22:00",
			ReservationRules: "Reservas de 1 hora, em horários inteiros, com múltiplas reservas por dia.",
			IsActive:         true,
		},
	}
}
