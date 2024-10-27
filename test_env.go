package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env:", err)
	} else {
		fmt.Println("Arquivo .env carregado com sucesso")
	}

	secret := os.Getenv("JWT_SECRET")
	fmt.Println("Valor de JWT_SECRET:", secret)
}
