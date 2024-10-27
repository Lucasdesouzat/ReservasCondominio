package services

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/Lucasdesouzat/ReservasCondominio/database"
	"github.com/Lucasdesouzat/ReservasCondominio/models"
)

// Segredo para o JWT (deve ser configurado no .env)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Inicialização para verificar o segredo JWT
func init() {
	log.Println("Segredo JWT carregado:", string(jwtSecret)) // Verifique se a chave está correta
}

// Função para registro de usuário
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Erro ao bindar JSON:", err) // Log detalhado
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Verificar os dados recebidos
	log.Printf("Dados recebidos: %+v\n", user)

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Erro ao gerar hash da senha:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar hash da senha"})
		return
	}
	user.Password = string(hashedPassword)

	// Converte o campo BirthDate para time.Time
	birthDate := user.BirthDate.ToTime()

	// Inserir usuário no banco de dados
	_, err = database.DB.NamedExec(`
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
		log.Println("Erro ao inserir usuário no banco de dados:", err) // Log detalhado
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar usuário"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário registrado com sucesso"})
}

// Função para autenticação de usuário
func LoginUser(c *gin.Context) {
	var user models.User
	var dbUser models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Erro ao bindar JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Procurar o usuário pelo email
	err := database.DB.Get(&dbUser, "SELECT * FROM users WHERE email = $1", user.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
		return
	}

	// Verificar a senha
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
		return
	}

	// Gerar o token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   dbUser.ID,
		"role": dbUser.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(), // Token expira em 72 horas
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
