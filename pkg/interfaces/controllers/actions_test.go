package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter(controller *ActionsController) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", controller.Ping)
	router.POST("/google-sheet", controller.GetGoogleSheetByID)
	router.POST("/notion", controller.GetNotion)
	return router
}

func TestPing(t *testing.T) {
	controller := NewActionsController(nil)
	router := setupTestRouter(controller)

	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"test":"test"`)
}
