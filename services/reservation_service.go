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

	// Regras para espaços específicos
	switch reservation.SpaceID {
	case 1: // Salão de Festas
		reservation.StartTime = time.Date(reservation.StartTime.Year(), reservation.StartTime.Month(), reservation.StartTime.Day(), 10, 0, 0, 0, reservation.StartTime.Location())
		reservation.EndTime = reservation.StartTime.Add(time.Hour * 14) // Das 10h até 00h (14 horas de duração)

		var existingReservation models.Reservation
		err := database.DB.Get(&existingReservation, `
            SELECT * FROM reservations 
            WHERE space_id = $1 AND status = $2 
            AND DATE(start_time) = DATE($3)`,
			reservation.SpaceID, constants.ReservationStatusConfirmed, reservation.StartTime)

		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "O Salão de Festas já está reservado para este dia."})
			return
		}

	case 2: // Brinquedoteca
		startHour := reservation.StartTime.Hour()
		if startHour < 8 || startHour >= 22 || reservation.StartTime.Minute() != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Horário inválido. A Brinquedoteca só permite reservas em horas inteiras das 8h às 22h."})
			return
		}
		reservation.EndTime = reservation.StartTime.Add(time.Hour) // 1 hora de duração

		var existingReservation models.Reservation
		err := database.DB.Get(&existingReservation, `
            SELECT * FROM reservations 
            WHERE space_id = $1 AND start_time = $2 
            AND status = $3`,
			reservation.SpaceID, reservation.StartTime, constants.ReservationStatusConfirmed)

		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "A Brinquedoteca já está reservada para esse horário."})
			return
		}

	case 3: // Academia
		if reservation.StartTime.Minute() != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Horário inválido. A Academia só permite reservas em horários inteiros."})
			return
		}
		reservation.EndTime = reservation.StartTime.Add(time.Hour) // 1 hora de duração

		var existingReservation models.Reservation
		err := database.DB.Get(&existingReservation, `
            SELECT * FROM reservations 
            WHERE space_id = $1 AND start_time = $2 
            AND status = $3`,
			reservation.SpaceID, reservation.StartTime, constants.ReservationStatusConfirmed)

		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "A Academia já está reservada para esse horário."})
			return
		}

	case 4, 5: // Quiosques
		startHour := reservation.StartTime.Hour()
		if !(startHour == 10 && reservation.EndTime.Hour() == 16) && !(startHour == 16 && reservation.EndTime.Hour() == 22) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Os Quiosques só permitem reservas para períodos: manhã (10h-16h) ou tarde (16h-22h)."})
			return
		}
		reservation.EndTime = reservation.StartTime.Add(time.Hour * 6)

	case 6: // Piscina
		startHour := reservation.StartTime.Hour()
		if startHour < 9 || startHour >= 22 || reservation.StartTime.Minute() != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Horário inválido. A Piscina permite reservas em horários inteiros das 9h às 22h."})
			return
		}
		reservation.EndTime = reservation.StartTime.Add(time.Hour) // 1 hora de duração
	}

	// Verifica se o usuário já possui uma reserva ativa
	var existingActiveReservation models.Reservation
	err := database.DB.Get(&existingActiveReservation, `
        SELECT * FROM reservations 
        WHERE user_id = $1 
        AND status = $2 AND end_time > NOW()`,
		reservation.UserID, constants.ReservationStatusConfirmed)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Você já possui uma reserva ativa."})
		return
	}

	// Configurar informações adicionais da reserva
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()
	reservation.Status = constants.ReservationStatusConfirmed

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

// Função para obter o histórico de reservas de um usuário
func GetReservationHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		log.Println("ID do usuário não encontrado no contexto")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário não encontrado"})
		return
	}

	var reservations []models.Reservation
	err := database.DB.Select(&reservations, `
        SELECT * FROM reservations 
        WHERE user_id = $1 
        ORDER BY start_time DESC
    `, userID)

	if err != nil {
		log.Println("Erro ao obter histórico de reservas:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter histórico de reservas"})
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
