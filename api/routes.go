package api

import (
	"github.com/Lucasdesouzat/ReservasCondominio/services"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Rota de autenticação
	router.POST("/register", services.RegisterUser)
	router.POST("/login", services.LoginUser)

	// Rota de teste para verificação de saúde do servidor
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Servidor está funcionando!"})
	})

	return router
}
