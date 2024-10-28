package models

import "time"

// Estrutura para representar uma reserva
type Reservation struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`       // Relacionado ao usuário que fez a reserva
	SpaceID   int       `json:"space_id" db:"space_id"`     // Relacionado ao espaço reservado
	StartTime time.Time `json:"start_time" db:"start_time"` // Hora de início da reserva
	EndTime   time.Time `json:"end_time" db:"end_time"`     // Hora de término da reserva
	Status    string    `json:"status" db:"status"`         // Status da reserva (Ex.: 'confirmado', 'pendente')
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
