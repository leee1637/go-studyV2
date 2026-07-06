package http_student

import (
	"net/http"
	"strconv"
	"study/internal/core/domain"

	"github.com/gin-gonic/gin"
)

func (s *StudentHandler) GetAll(g *gin.Context) {

	page, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(g.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 20
	}

	result, err := s.StudentService.GetAll(g.Request.Context(), page, limit)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	g.JSON(http.StatusOK, result)
}

func (s *StudentHandler) GetByID(g *gin.Context) {
	requestID, exists := g.Get("userID")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Нету авторизации"})
		return
	}

	role, _ := g.Get("role")

	studentID, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
		return
	}

	result, err := s.StudentService.GetByIDStudent(g.Request.Context(), requestID.(int), domain.Role(role.(string)), studentID)
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, result)

}

func (s *StudentHandler) GetByGroup(g *gin.Context) {
	requestID, exists := g.Get("userID")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Нету авторизации"})
		return
	}

	role, _ := g.Get("role")

	group := g.Param("group")
	if group == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "ошибка парса параметра group"})
		return
	}

	students, err := s.StudentService.GetByGroup(g.Request.Context(), requestID.(int), domain.Role(role.(string)), string(group))
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, students)
}

func (s *StudentHandler) UpdateStudent(g *gin.Context) {
	requestID, exists := g.Get("userID")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Нету авторизации"})
		return
	}

	role, _ := g.Get("role")

	studentID, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
		return
	}

	var input domain.UpdateStudentDTO
	if err := g.ShouldBindJSON(&input); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	err = s.StudentService.UpdateStudent(g.Request.Context(), requestID.(int), domain.Role(role.(string)), studentID, input)
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"result": "успешно обновлён"})

}

func (s *StudentHandler) DeleteStudent(g *gin.Context) {
	requestID, exists := g.Get("userID")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Нету авторизации"})
		return
	}

	role, _ := g.Get("role")

	studentID, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
		return
	}

	err = s.StudentService.DeleteStudent(g.Request.Context(), requestID.(int), domain.Role(role.(string)), studentID)
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"result": "успешно удалён"})
}
