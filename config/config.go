package config

import (
	"log"
	// Biblioteca para carregamento de arquivos `.env`
	"github.com/joho/godotenv"
)

func LoadConfig() {
	// Carrega variáveis do arquivo `.env`
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar arquivo .env")
	}
}
