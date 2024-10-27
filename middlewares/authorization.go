package middlewares

import (
	"net/http"

	"github.com/Lucasdesouzat/ReservasCondominio/pkg/constants"
	"github.com/gin-gonic/gin"
)

// Middleware para verificar permissões com base no papel
func AuthorizeRole(requiredRole constants.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != string(requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada"})
			c.Abort()
			return
		}
		c.Next()
	}
}
