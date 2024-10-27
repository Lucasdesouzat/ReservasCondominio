package services

import (
	"log"
	"net/http"

	"github.com/Lucasdesouzat/ReservasCondominio/database"
	"github.com/Lucasdesouzat/ReservasCondominio/models"
	"github.com/gin-gonic/gin"
)

// Função para criar um novo espaço
func CreateSpace(c *gin.Context) {
	var space models.Space
	if err := c.ShouldBindJSON(&space); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Inserir espaço no banco de dados
	_, err := database.DB.NamedExec(`
        INSERT INTO spaces (name, description, max_reservations, max_occupancy, price, available_from, available_until, amenities, requires_approval, reservation_rules, is_active, image_url) 
        VALUES (:name, :description, :max_reservations, :max_occupancy, :price, :available_from, :available_until, :amenities, :requires_approval, :reservation_rules, :is_active, :image_url)
    `, &space)
	if err != nil {
		log.Println("Erro ao criar espaço:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar espaço"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Espaço criado com sucesso"})
}

// Função para obter todos os espaços
func GetAllSpaces(c *gin.Context) {
	var spaces []models.Space
	err := database.DB.Select(&spaces, "SELECT * FROM spaces")
	if err != nil {
		log.Println("Erro ao obter espaços:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter espaços"})
		return
	}

	c.JSON(http.StatusOK, spaces)
}
