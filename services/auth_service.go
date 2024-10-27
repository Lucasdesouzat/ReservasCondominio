package services

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/Lucasdesouzat/ReservasCondominio/database"
	"github.com/Lucasdesouzat/ReservasCondominio/models"
	// Incluindo utils para funções de validação
)

// Função que realiza o registro do usuário no banco de dados
func RegisterUserService(user models.User) error {
	// Converte o campo BirthDate para time.Time
	birthDate := user.BirthDate.ToTime()

	// Inserir usuário no banco de dados
	_, err := database.DB.NamedExec(`
		INSERT INTO users (
			first_name, 
			last_name, 
			cpf, 
			birth_date, 
			profile_picture, 
			phone_1, 
			phone_2, 
			email, 
			password, 
			role, 
			status, 
			auth_method
		) VALUES (
			:first_name, 
			:last_name, 
			:cpf, 
			:birth_date, 
			:profile_picture, 
			:phone_1, 
			:phone_2, 
			:email, 
			:password, 
			:role, 
			:status, 
			:auth_method
		)`, map[string]interface{}{
		"first_name":      user.FirstName,
		"last_name":       user.LastName,
		"cpf":             user.CPF,
		"birth_date":      birthDate,
		"profile_picture": user.ProfilePicture,
		"phone_1":         user.Phone1,
		"phone_2":         user.Phone2,
		"email":           user.Email,
		"password":        user.Password,
		"role":            user.Role,
		"status":          user.Status,
		"auth_method":     user.AuthMethod,
	})
	if err != nil {
		log.Println("Erro ao inserir usuário no banco de dados:", err)
		return err
	}

	return nil
}

// Função que realiza o login do usuário e gera o token JWT
func LoginService(email, password string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Erro: JWT_SECRET não está configurado corretamente")
	}

	var dbUser models.User
	err := database.DB.Get(&dbUser, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		log.Println("Erro ao buscar usuário no banco de dados:", err)
		return "", errors.New("Email ou senha inválidos")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		log.Println("Erro na comparação de senha:", err)
		return "", errors.New("Email ou senha inválidos")
	}

	// Gerar o token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   dbUser.ID,
		"role": dbUser.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(), // Token expira em 72 horas
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("Erro ao gerar token JWT:", err)
		return "", errors.New("Erro ao gerar token")
	}

	return tokenString, nil
}
