package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Driver para PostgreSQL
)

var DB *sqlx.DB

func Connect() {
	// Obtém a URL do banco de dados a partir do arquivo .env
	databaseURL := os.Getenv("DATABASE_URL")

	// Conecta ao banco de dados PostgreSQL
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	// Testa a conexão com o banco de dados
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao testar conexão com o banco de dados:", err)
	}

	DB = db
	log.Println("Conexão ao banco de dados estabelecida com sucesso!")
}
