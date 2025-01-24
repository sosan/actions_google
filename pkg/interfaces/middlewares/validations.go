package middlewares

import (
	"actions_google/pkg/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateOnGetWorkflow() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: validations
		ctx.Next()
	}
}

func ValidateUserAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func ValidateGetGoogleSheet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var currentReq models.ActionsCommand
		if err := ctx.ShouldBindBodyWithJSON(&currentReq); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}

		if !models.ValidCommandTypes[*currentReq.Type] {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}

		if !models.ValidActionsTypes[currentReq.Actions.Type] {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}

		if currentReq.Actions.RedirectURL == "" {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}

		// TODO:
		// check operation // getallcontent
		// credentialid required and not null
		// sub required and not null
		// workflowid required and not null

		ctx.Set(models.ActionGoogleKey, currentReq)
		ctx.Next()
	}
}

func ValidateNotionFields() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var currentReq models.ActionsCommand
		if err := ctx.ShouldBindBodyWithJSON(&currentReq); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}

		if !models.ValidCommandTypes[*currentReq.Type] {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}

		if !models.ValidActionsTypes[currentReq.Actions.Type] {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}

		// TODO:
		// check operation // getallcontent
		// credentialid required and not null
		// sub required and not null
		// workflowid required and not null

		ctx.Set(models.ActionNotionKey, currentReq)
		ctx.Next()
	}
}
