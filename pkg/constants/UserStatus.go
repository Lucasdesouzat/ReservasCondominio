package constants

type UserStatus string

// Constantes para status de usuários
const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBanned   UserStatus = "banned"
)
