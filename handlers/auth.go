package handlers

import (
	"log"
	"net/http"

	"github.com/Lucasdesouzat/ReservasCondominio/models"
	"github.com/Lucasdesouzat/ReservasCondominio/services"
	"github.com/Lucasdesouzat/ReservasCondominio/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Estrutura para dados de login recebidos
type DadosLogin struct {
	Email string `json:"email"`
	Senha string `json:"password"`
}

// Handler para registrar um novo usuário
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Erro ao bindar JSON para registro:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de registro inválidos"})
		return
	}

	log.Printf("Tentativa de registro com dados: %+v\n", user)

	// Validação do CPF
	if !utils.ValidarCPF(user.CPF) {
		log.Println("CPF inválido:", user.CPF)
		c.JSON(http.StatusBadRequest, gin.H{"error": "CPF inválido"})
		return
	}

	// Validação do Email
	if !utils.ValidarFormatoEmail(user.Email) {
		log.Println("Email inválido:", user.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email inválido"})
		return
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Erro ao gerar hash da senha:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar hash da senha"})
		return
	}
	user.Password = string(hashedPassword)

	// Registrar o usuário no banco de dados usando o serviço
	err = services.RegisterUserService(user)
	if err != nil {
		if err.Error() == "usuário já registrado com esse CPF" {
			c.JSON(http.StatusConflict, gin.H{"error": "Usuário já registrado com este CPF"})
		} else {
			log.Println("Erro ao registrar usuário:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar usuário"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário registrado com sucesso"})
}

// Handler de login usando o contexto do Gin
func LoginHandler(c *gin.Context) {
	var dadosLogin DadosLogin
	if err := c.ShouldBindJSON(&dadosLogin); err != nil {
		log.Println("Erro ao bindar JSON para login:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar requisição"})
		return
	}

	log.Println("Tentativa de login para o email:", dadosLogin.Email)

	// Chamar o serviço de autenticação e receber o token
	token, err := services.LoginService(dadosLogin.Email, dadosLogin.Senha)
	if err != nil {
		log.Println("Erro na autenticação:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login realizado com sucesso", "token": token})
}
