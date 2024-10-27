package api

import (
	"github.com/Lucasdesouzat/ReservasCondominio/middlewares"
	"github.com/Lucasdesouzat/ReservasCondominio/pkg/constants"
	"github.com/Lucasdesouzat/ReservasCondominio/services"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Rotas públicas (Registro e Login)
	router.POST("/register", services.RegisterUser)
	router.POST("/login", services.LoginUser) // <--- Verifique se esta linha está presente

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

	// Exemplos de rotas para reservas
	auth.POST("/reservations", middlewares.AuthorizeRole(constants.RoleResident), services.CreateReservation)
	auth.GET("/reservations", middlewares.AuthorizeRole(constants.RoleResident), services.GetUserReservations)

	return router
}
