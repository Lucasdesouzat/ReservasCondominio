package api

import (
	"github.com/Lucasdesouzat/ReservasCondominio/handlers"
	"github.com/Lucasdesouzat/ReservasCondominio/middlewares"
	"github.com/Lucasdesouzat/ReservasCondominio/pkg/constants"
	"github.com/Lucasdesouzat/ReservasCondominio/services"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Rotas públicas (Registro e Login)
	router.POST("/register", handlers.RegisterUser) // Usando handlers.RegisterUser
	router.POST("/login", handlers.LoginHandler)    // Usando handlers.LoginHandler

	// Rota de teste para verificação de saúde do servidor
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Servidor está funcionando!"})
	})

	// Rotas protegidas com autenticação JWT
	auth := router.Group("/api")
	auth.Use(middlewares.JWTAuthMiddleware()) // Middleware para autenticação JWT

	// Rotas de administração de espaços (somente para administradores)
	auth.POST("/spaces", middlewares.AuthorizeRole(constants.RoleAdmin), services.CreateSpace)
	auth.GET("/spaces", middlewares.AuthorizeRole(constants.RoleAdmin), services.GetAllSpaces)

	// Exemplos de rotas para reservas (permitindo admin e resident)
	auth.POST("/reservations", middlewares.AuthorizeRole(constants.RoleAdmin, constants.RoleResident, constants.RoleOwner), services.CreateReservation)
	auth.GET("/reservations", middlewares.AuthorizeRole(constants.RoleAdmin, constants.RoleResident, constants.RoleOwner), services.GetUserReservations)

	// Rota para cancelar reservas (permitindo admin e owner)
	auth.POST("/reservations/cancel", middlewares.AuthorizeRole(constants.RoleAdmin, constants.RoleResident, constants.RoleOwner), services.CancelReservation)

	return router
}
