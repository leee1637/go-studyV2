package transport_teacher

import (
	"net/http"
	"strconv"
	"study/internal/core/domain"

	"github.com/gin-gonic/gin"
)

func (t *TeacherHandler) GetAll(g *gin.Context) {
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

	result, err := t.TeacherService.GetAll(g.Request.Context(), page, limit)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	g.JSON(http.StatusOK, result)
}

func (t *TeacherHandler) GetByID(g *gin.Context) {
	idReq, exists := g.Get("userID")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизоан"})
		return
	}

	role, _ := g.Get("role")

	teacherID, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
		return
	}

	result, err := t.TeacherService.GetByID(g.Request.Context(), idReq.(int), teacherID, domain.Role(role.(string)))
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, result)
}

func (t *TeacherHandler) UpdateTeacher(g *gin.Context) {
	idReq, exists := g.Get("userID")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизоан"})
		return
	}

	role, _ := g.Get("role")

	teacherID, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
		return
	}

	var input domain.UpdateTeachertDTO
	if err := g.ShouldBindJSON(&input); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	err = t.TeacherService.UpdateTeacher(g.Request.Context(), idReq.(int), teacherID, domain.Role(role.(string)), input)
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"status": "успешно обновлен пользователь"})
}

func (t *TeacherHandler) DeleteTeacher(g *gin.Context) {
	idReq, exists := g.Get("userID")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизоан"})
		return
	}

	role, _ := g.Get("role")

	teacherID, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
		return
	}

	err = t.TeacherService.DeleteTeacher(g.Request.Context(), idReq.(int), teacherID, domain.Role(role.(string)))
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"status": "успешно удалён пользователь"})
}

func (t *TeacherHandler) GetByGroup(g *gin.Context) {
	idReq, exists := g.Get("userID")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизоан"})
		return
	}

	role, _ := g.Get("role")

	teacherGroup := g.Param("group")
	if teacherGroup == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Не задан group"})
		return
	}

	tec, err := t.TeacherService.GetByGroup(g.Request.Context(), idReq.(int), teacherGroup, domain.Role(role.(string)))
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, tec)
}

func (t *TeacherHandler) AddGroup(g *gin.Context) {
	_, exists := g.Get("userID")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизоан"})
		return
	}

	role, _ := g.Get("role")

	teacherID, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
		return
	}

	teacherGroup := g.Param("group")
	if teacherGroup == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Не задан group"})
		return
	}

	err = t.TeacherService.AddGroup(g.Request.Context(), teacherID, teacherGroup, domain.Role(role.(string)))
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"status": "Группа добалвена"})
}

func (t *TeacherHandler) DeleteGroup(g *gin.Context) {
	_, exists := g.Get("userID")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизоан"})
		return
	}

	role, _ := g.Get("role")

	teacherGroup := g.Param("group")
	if teacherGroup == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Не задан group"})
		return
	}

	err := t.TeacherService.DeleteGroup(g.Request.Context(), teacherGroup, domain.Role(role.(string)))
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"status": "Группа удалена"})
}
