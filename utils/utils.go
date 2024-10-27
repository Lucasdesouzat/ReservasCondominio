package utils

import (
	"log"
	"regexp"
	"strings"
)

// Função para validar o formato do email usando Regex
func ValidarFormatoEmail(email string) bool {
	const regexEmail = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regexEmail)
	return re.MatchString(email)
}

// Função de validação simples para CPF
func ValidarCPF(cpf string) bool {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")

	// Verifica se o CPF tem 11 dígitos
	if len(cpf) != 11 {
		log.Println("Erro: CPF deve conter 11 dígitos")
		return false
	}

	// Implementar aqui a lógica de validação do CPF
	// Por exemplo, verificar dígitos verificadores

	return true
}
