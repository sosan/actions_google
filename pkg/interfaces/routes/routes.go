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
		// api.GET("/ping", common.Ping)

		actions := api.Group("/actions")
		{
			actions.GET("/google", dependencies.ActionsController.Ping)
			actions.POST("/google/sheets", middlewares.ValidateGetGoogleSheet(), dependencies.ActionsController.GetGoogleSheetByID)
		}
	}
}

func ErrRouter(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Page not found",
	})
}
