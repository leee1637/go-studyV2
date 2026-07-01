package http_student

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Input struct {
}

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

	result, err := s.StudentService.GetAll(g, page, limit)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	g.JSON(http.StatusOK, result)
}
