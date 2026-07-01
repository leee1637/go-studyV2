package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(g *gin.Context) {
		au := g.GetHeader("Authorization")

		if au == "" {
			g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует заголовок Authorization"})
			return
		}

		parts := strings.Split(au, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат заголовка Authorization"})
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Невалидный или протухший токен"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Не удалось прочитать данные токена"})
			return
		}

		g.Set("userID", int(claims["user_id"].(float64)))
		g.Set("role", claims["role"].(string))

		g.Next()
	}
}

func RequireRole(allowedRole ...string) gin.HandlerFunc {
	return func(g *gin.Context) {
		userRole, exists := g.Get("role")

		if !exists {
			g.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Роль не определена"})
			return
		}

		for _, v := range allowedRole {
			if userRole.(string) == v {
				g.Next()
				return
			}
		}

		g.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "У вас нет прав для этого действия"})

	}
}
