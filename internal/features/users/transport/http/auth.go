package http_transport

import (
	"net/http"
	"study/internal/core/domain"

	"github.com/gin-gonic/gin"
)

type SignUpInput struct {
	Login    string      `json:"login" binding:"required"`
	Password string      `json:"password" binding:"required"`
	Role     domain.Role `json:"role" binding:"required"`
}

type SignInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *AuthHandler) SignUp(g *gin.Context) {
	var input SignUpInput

	err := g.ShouldBindJSON(&input)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных" + err.Error()})
		return
	}

	err = a.authService.SignUp(g.Request.Context(), input.Login, input.Password, input.Role)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

func (a *AuthHandler) SignIn(g *gin.Context) {
	var input SignInInput

	err := g.ShouldBindJSON(&input)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных" + err.Error()})
		return
	}

	token, err := a.authService.SignIn(g.Request.Context(), input.Login, input.Password)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"token": token})

}
