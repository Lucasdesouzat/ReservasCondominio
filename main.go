package main

import (
	"log"
	"os"

	"github.com/Lucasdesouzat/ReservasCondominio/database"
	"github.com/Lucasdesouzat/ReservasCondominio/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Para carregar o .env
)

func main() {
	// Log para verificar o diretório atual
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Erro ao obter o diretório atual:", err)
	}
	log.Println("Diretório atual:", dir)

	// Tentando carregar o arquivo .env
	log.Println("Tentando carregar o arquivo .env...")

	// Carregar variáveis de ambiente do arquivo .env
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env:", err)
	} else {
		log.Println("Arquivo .env carregado com sucesso")
	}

	// Verifique se o JWT_SECRET está configurado
	secret := os.Getenv("JWT_SECRET")
	log.Println("Valor lido de JWT_SECRET:", secret) // Adiciona log para depurar o valor
	if secret == "" {
		log.Fatal("Erro: JWT_SECRET não está configurado corretamente")
	} else {
		log.Println("JWT_SECRET carregado com sucesso:", secret) // Log para verificar o carregamento
	}

	// Inicializar conexão com o banco de dados
	database.Connect()

	// Configurar rotas do servidor Gin
	r := gin.Default()

	// Definir rotas
	r.POST("/register", handlers.RegisterUser)
	r.POST("/login", handlers.LoginUser)

	// Iniciar o servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
