package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Lucasdesouzat/ReservasCondominio/services" // Import correto do serviço de autenticação
)

// Estrutura para dados de login recebidos
type DadosLogin struct {
	Email string `json:"email"`
	Senha string `json:"password"`
}

// Handler de login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse do corpo da requisição
	var dadosLogin DadosLogin
	err := json.NewDecoder(r.Body).Decode(&dadosLogin)
	if err != nil {
		log.Println("Erro ao decodificar o JSON:", err)
		http.Error(w, "Erro ao processar requisição", http.StatusBadRequest)
		return
	}

	// Log das informações recebidas
	log.Println("Iniciando login para o email:", dadosLogin.Email)

	// Chamar o serviço de autenticação
	err = services.LoginService(dadosLogin.Email, dadosLogin.Senha)
	if err != nil {
		log.Println("Erro no login:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Resposta de sucesso
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login realizado com sucesso"))
}
