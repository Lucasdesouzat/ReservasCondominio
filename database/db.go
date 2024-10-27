package database

import (
	"log"

	"github.com/Lucasdesouzat/ReservasCondominio/models" // Import do modelo de usuário
)

// Função para buscar usuário pelo email no banco de dados
func BuscarUsuarioPorEmail(email string) (*models.User, error) {
	var usuario models.User
	err := DB.QueryRow("SELECT id, email, password FROM usuarios WHERE email_nm = $1", email).Scan(&usuario.ID, &usuario.Email, &usuario.Password)

	if err != nil {
		log.Println("Erro ao buscar usuário no banco de dados:", err)
		return nil, err
	}

	return &usuario, nil
}
