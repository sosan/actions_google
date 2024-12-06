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

		// workflows := api.Group("/workflows")
		// {
		// 	workflows.GET("/:iduser/workflow/:idworkflow/:usertoken", dependencies.AuthController.VerifyUserTokenForMiddleware, middlewares.ValidateOnGetWorkflow(), dependencies.WorkflowController.GetWorkflow)
		// 	workflows.GET("/:iduser/:usertoken", dependencies.AuthController.VerifyUserTokenForMiddleware, middlewares.ValidateOnGetWorkflow(), dependencies.WorkflowController.GetAllWorkflows)
		// 	workflows.POST("", middlewares.ValidateOnCreateWorkflow(), dependencies.WorkflowController.CreateWorkflow)
		// 	workflows.PUT("/:id", middlewares.ValidateOnUpdateWorkflow(), dependencies.WorkflowController.UpdateWorkflow)
		// }

		// users := api.Group("/users")
		// {
		// 	users.POST("", middlewares.ValidateUser(), dependencies.UserController.SyncUseWrithIDProvider)
		// 	users.GET("/:stub", dependencies.UserController.GetUserByStub)
		// }

		// credentials := api.Group("/credentials")
		// {
		// 	// credentials.POST("", middlewares.ValidateUser(), userController.SyncUseWrithIDProvider)
		// 	credentials.GET("/:iduser/:usertoken", dependencies.AuthController.VerifyUserTokenForMiddleware, dependencies.CredentialController.GetAllCredentials)
		// }

		// dashboard := api.Group("/dashboard")
		// {
		// 	dashboard.GET("/:iduser", middlewares.ValidateUserAuth(), dependencies.DashboardController.GetUserDashboardByID)
		// }

		// auth := api.Group("/auth")
		// {
		// 	auth.GET("/verify/:usertoken", dependencies.AuthController.VerifyUserToken)
		// }

		// credentialsGoogle := api.Group("/google")
		// {
		// 	credentialsGoogle.POST("/credential", middlewares.ValidateOnCreateCredential(), dependencies.CredentialController.CreateCredential)
		// 	credentialsGoogle.POST("/exchange", middlewares.ValidateOnExchangeCredential(), dependencies.CredentialController.ExchangeGoogleCode)
		// }

		actions := api.Group("/actions")
		{
			actions.POST("/google/sheets", middlewares.ValidateGetGoogleSheet(), dependencies.ActionsController.GetGoogleSheetByID)
		}
	}
}

func ErrRouter(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Page not found",
	})
}


