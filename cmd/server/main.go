package main

import (
	"log"
	"os"

	"github.com/Lucasdesouzat/ReservasCondominio/api"
	"github.com/Lucasdesouzat/ReservasCondominio/database"
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

	// Configurar rotas do servidor usando api.SetupRouter()
	router := api.SetupRouter()

	// Iniciar o servidor usando a variável 'router'
	log.Println("Servidor rodando na porta 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
