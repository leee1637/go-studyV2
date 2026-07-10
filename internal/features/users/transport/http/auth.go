package http_transport

import (
	"io"
	"net/http"
	"study/internal/core/domain"
	"time"

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

func (a *AuthHandler) UploadCSV(g *gin.Context) {
	fileHeader, err := g.FormFile("file")
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен"})
		return
	}

	const MaxSize = 5 * 1024 * 1024
	if fileHeader.Size > MaxSize {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Файл слишком большой (макс. 5 МБ)"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения файла"})
		return
	}

	result, err := a.regService.ProcessCSV(g.Request.Context(), fileBytes)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.Processed == 0 && len(result.Warnings) > 0 {
		g.JSON(http.StatusUnprocessableEntity, result)
		return
	}

	g.JSON(http.StatusOK, result)
}

func (a *AuthHandler) VerifyCSVRegistration(g *gin.Context) {
	tokenStr := g.Param("token")
	token, err := uuid.Parse(tokenStr)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена"})
		return
	}

	regReq, err := a.regService.GetByToken(g.Request.Context(), token)
	if err != nil {
		g.JSON(http.StatusGone, gin.H{"error": "Ссылка недействительна"})
		return
	}

	if regReq.Status != "pending" {
		g.JSON(http.StatusGone, gin.H{"error": "Заявка уже обработана"})
		return
	}

	if time.Now().After(regReq.ExpiresAt) {
		a.regService.DeleteByToken(g.Request.Context(), token)
		g.JSON(http.StatusGone, gin.H{"error": "Срок действия ссылки истёк"})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"token": token.String(),
		"fio":   regReq.FIO,
		"email": regReq.Email,
		"role":  regReq.Role,
	})
}

func (a *AuthHandler) CompleteCSVRegistration(g *gin.Context) {
	tokenStr := g.Param("token")
	token, err := uuid.Parse(tokenStr)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена"})
		return
	}

	var input struct {
		Password string `json:"password" binding:"required"`
	}
	if err := g.ShouldBindJSON(&input); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "password обязателен"})
		return
	}

	err = a.regService.CompleteRegistration(g.Request.Context(), token, input.Password)
	if err != nil {
		switch {
		case err.Error() == "Заявка не найдена" || err.Error() == "Срок действия ссылки истёк":
			g.JSON(http.StatusGone, gin.H{"error": err.Error()})
		case err.Error() == "Пользователь с таким email уже существует":
			g.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	g.JSON(http.StatusOK, gin.H{"message": "Регистрация завершена! Теперь вы можете войти."})
}
