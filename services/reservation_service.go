package services

import (
	"log"
	"net/http"
	"time"

	"github.com/Lucasdesouzat/ReservasCondominio/database"
	"github.com/Lucasdesouzat/ReservasCondominio/models"
	"github.com/Lucasdesouzat/ReservasCondominio/pkg/constants"
	"github.com/gin-gonic/gin"
)

// Função para criar uma nova reserva
func CreateReservation(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		log.Println("Erro ao bindar JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Verificar se o horário de início é válido
	startHour := reservation.StartTime.Hour()
	startMinute := reservation.StartTime.Minute()
	if startHour < 0 || startHour > 23 || startMinute != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Horário inválido. Os horários disponíveis são das 10h às 14h, e devem ser em horários fixos (ex.: 10:00, 11:00)."})
		return
	}

	// Definir o horário de término como 1 hora após o início
	reservation.EndTime = reservation.StartTime.Add(time.Hour)

	// Definindo os valores de CreatedAt e UpdatedAt e o status
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()
	reservation.Status = constants.ReservationStatusConfirmed

	// Log dos dados da reserva recebidos
	log.Printf("Dados da reserva recebidos: %+v\n", reservation)

	// Verifica se o usuário já possui uma reserva ativa para o mesmo espaço
	var existingActiveReservation models.Reservation
	err := database.DB.Get(&existingActiveReservation, `
		SELECT * FROM reservations 
		WHERE user_id = $1 AND space_id = $2 
		AND status = $3 AND end_time > NOW()`,
		reservation.UserID, reservation.SpaceID, constants.ReservationStatusConfirmed)

	if err == nil {
		log.Printf("Reserva ativa encontrada: %+v\n", existingActiveReservation)
		c.JSON(http.StatusConflict, gin.H{"error": "Você já possui uma reserva ativa para este espaço."})
		return
	} else {
		log.Println("Nenhuma reserva ativa encontrada ou erro:", err)
	}

	// Verifica se já existe uma reserva no mesmo horário (considerando apenas reservas confirmadas)
	var existingReservation models.Reservation
	err = database.DB.Get(&existingReservation, `
		SELECT * FROM reservations 
		WHERE space_id = $1 
		AND start_time = $2 
		AND status = $3`,
		reservation.SpaceID, reservation.StartTime, constants.ReservationStatusConfirmed)

	if err == nil {
		log.Printf("Já existe uma reserva para o horário: %+v\n", existingReservation)
		c.JSON(http.StatusConflict, gin.H{"error": "Já existe uma reserva para a academia nesse horário."})
		return
	} else {
		log.Println("Nenhuma reserva existente encontrada ou erro:", err)
	}

	// Inserir a reserva no banco de dados
	_, err = database.DB.NamedExec(`
		INSERT INTO reservations (user_id, space_id, start_time, end_time, status, created_at, updated_at) 
		VALUES (:user_id, :space_id, :start_time, :end_time, :status, :created_at, :updated_at)
	`, &reservation)
	if err != nil {
		log.Println("Erro ao criar reserva:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar reserva"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Reserva criada com sucesso"})
}

// Função para obter as reservas de um usuário
func GetUserReservations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		log.Println("ID do usuário não encontrado no contexto")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário não encontrado"})
		return
	}

	var reservations []models.Reservation
	err := database.DB.Select(&reservations, "SELECT * FROM reservations WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Erro ao obter reservas:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter reservas"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

// Função para cancelar uma reserva
func CancelReservation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		log.Println("ID do usuário não encontrado no contexto")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário não encontrado"})
		return
	}

	role := c.GetHeader("role")
	isAdmin := (constants.UserRole(role) == constants.RoleAdmin)

	var reservation models.Reservation

	if isAdmin {
		var requestBody struct {
			CPF     string `json:"cpf"`
			SpaceID int    `json:"space_id"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
			return
		}

		err := database.DB.Get(&reservation, `
			SELECT * FROM reservations 
			WHERE user_id = (SELECT id FROM users WHERE cpf = $1) AND space_id = $2 
			AND status = $3 AND start_time > NOW() 
			ORDER BY created_at DESC LIMIT 1`, requestBody.CPF, requestBody.SpaceID, constants.ReservationStatusConfirmed)

		if err != nil {
			log.Println("Erro ao encontrar a reserva:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Nenhuma reserva encontrada para esse usuário e espaço."})
			return
		}
	} else {
		var requestBody struct {
			SpaceID int `json:"space_id"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
			return
		}

		err := database.DB.Get(&reservation, `
			SELECT * FROM reservations 
			WHERE user_id = $1 AND space_id = $2 AND status = $3 AND start_time > NOW() 
			ORDER BY created_at DESC LIMIT 1`, userID, requestBody.SpaceID, constants.ReservationStatusConfirmed)

		if err != nil {
			log.Println("Erro ao encontrar a reserva:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Nenhuma reserva encontrada ou não pode ser cancelada."})
			return
		}
	}

	// Permitir cancelamento
	reservation.Status = constants.ReservationStatusCancelled

	// Atualizar a reserva no banco de dados
	_, err := database.DB.NamedExec(`
		UPDATE reservations SET status = :status, updated_at = :updated_at WHERE id = :id
	`, map[string]interface{}{
		"status":     reservation.Status,
		"updated_at": time.Now(),
		"id":         reservation.ID,
	})

	if err != nil {
		log.Println("Erro ao cancelar reserva:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao cancelar reserva"})
		return
	}

	// Log para liberar o horário
	log.Printf("Horário liberado para o espaço %d para o usuário %d\n", reservation.SpaceID, userID)

	c.JSON(http.StatusOK, gin.H{"message": "Reserva cancelada com sucesso"})
}
