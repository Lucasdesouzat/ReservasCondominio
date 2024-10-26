package main

import (
	"github.com/Lucasdesouzat/ReservasCondominio/api"
	"github.com/Lucasdesouzat/ReservasCondominio/config"
	"github.com/Lucasdesouzat/ReservasCondominio/database"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Carrega variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	// Verifica se a variável JWT_SECRET foi carregada
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("Erro: a variável JWT_SECRET não está definida no arquivo .env")
	} else {
		log.Println("JWT_SECRET carregado com sucesso:", jwtSecret)
	}

	// Testa se a variável DATABASE_URL foi carregada
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL não definida no arquivo .env")
	} else {
		log.Println("Conectando ao banco de dados com a URL:", dbURL)
	}

	// Carrega configurações adicionais
	config.LoadConfig()

	// Conecta ao banco de dados
	database.Connect()

	// Configura as rotas da API
	router := api.SetupRouter()

	// Inicia o servidor
	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
