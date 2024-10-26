package constants

// Definição dos tipos de papéis para o usuário
type UserRole string

// Constantes para papéis de usuários
const (
	RoleResident UserRole = "resident"
	RoleOwner    UserRole = "owner"
	RoleAdmin    UserRole = "admin"
)
