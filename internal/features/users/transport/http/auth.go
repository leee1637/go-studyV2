package http_transport

import (
	"net/http"
	"study/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SignUpInput struct {
	Email       string      `json:"email" binding:"required"`
	Password    string      `json:"password" binding:"required"`
	Role        domain.Role `json:"role" binding:"required"`
	FIO         string      `json:"fio" binding:"required"`
	GroupName   []string    `json:"group_name"`
	PhoneNumber *string     `json:"phone_number"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshInput struct {
	RefreshInput string `json:"refresh_token"`
}

func (a *AuthHandler) SignUp(g *gin.Context) {
	var input SignUpInput

	err := g.ShouldBindJSON(&input)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных" + err.Error()})
		return
	}

	dto := domain.SignUpDTO{
		Email:       input.Email,
		Password:    input.Password,
		Role:        input.Role,
		FIO:         input.FIO,
		GroupName:   input.GroupName,
		PhoneNumber: input.PhoneNumber,
	}

	err = a.authService.SignUp(g.Request.Context(), dto)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	fastToken, refreshToken, err := a.authService.SignIn(g.Request.Context(), input.Email, input.Password)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"fast_token": fastToken,
		"refresh_token": refreshToken})

}

func (a *AuthHandler) Refresh(g *gin.Context) {
	var r RefreshInput

	err := g.ShouldBindJSON(&r)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Нету refresh_token - обязательно"})
		return
	}

	fastToken, refreshToken, err := a.authService.Refresh(g.Request.Context(), r.RefreshInput)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка парса токенов"})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"access_token":  fastToken,
		"refresh_token": refreshToken,
	})
}

func (a *AuthHandler) Logout(g *gin.Context) {
	var input RefreshInput
	if err := g.ShouldBindJSON(&input); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token обязателен"})
		return
	}

	err := a.authService.Logout(g.Request.Context(), input.RefreshInput)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"message": "Успешный выход"})
}

func (a *AuthHandler) VerifyEmail(g *gin.Context) {

	tokenStr := g.Param("token")
	token, err := uuid.Parse(tokenStr)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена"})
		return
	}

	err = a.authService.ConfirmEmail(g.Request.Context(), token)
	if err != nil {
		g.JSON(http.StatusGone, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"message": "Email успешно подтверждён! Теперь вы можете войти.",
	})
}
