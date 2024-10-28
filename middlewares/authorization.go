package middlewares

import (
	"log"
	"net/http"

	"github.com/Lucasdesouzat/ReservasCondominio/pkg/constants"
	"github.com/gin-gonic/gin"
)

// Middleware para verificar permissões com base em múltiplos papéis de tipo UserRole
func AuthorizeRole(requiredRoles ...constants.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada"})
			c.Abort()
			return
		}

		// Converter a role para string para comparação
		roleStr := role.(string)

		// Verificar se o papel do usuário está entre os papéis permitidos
		for _, requiredRole := range requiredRoles {
			if roleStr == string(requiredRole) {
				c.Next()
				return
			}
		}

		log.Printf("Permissão negada: papel atual %v, papéis permitidos %v\n", roleStr, requiredRoles)
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada"})
		c.Abort()
	}
}
