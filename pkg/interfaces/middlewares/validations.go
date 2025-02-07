package middlewares

import (
	"actions_google/pkg/domain/models"
	"net/http"
	"strings"

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
		// required sub ?
		if strings.TrimSpace(currentReq.Actions.Sub) == "" {
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
		// check spreadsheetID is valid
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
		// Frontend did first validation, backend second validation and third validation this one
		if !strings.Contains(currentReq.Actions.Document, models.NotionHost) {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}
		// Frontend did first validation, backend second validation and third validation this one
		if strings.TrimSpace(currentReq.Actions.Document) == "" {
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
