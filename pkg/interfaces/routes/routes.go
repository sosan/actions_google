package routes

import (
	"actions_google/pkg/dimodel"
	"actions_google/pkg/interfaces/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine, dependencies *dimodel.Dependencies) {
	app.NoRoute(ErrRouter)

	// Routes in groups
	api := app.Group("/api")
	{
		actions := api.Group("/actions")
		{
			actions.GET("/ping", dependencies.ActionsController.Ping)
			actions.POST("/google/sheets", middlewares.ValidateGetGoogleSheet(), dependencies.ActionsController.GetGoogleSheetByID)
			actions.POST("/notion", middlewares.ValidateNotionFields(), dependencies.ActionsController.GetNotion)
		}
	}
}

func ErrRouter(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Page not found",
	})
}
