package middlewares

import (
	"actions_google/pkg/domain/models"
	"log"
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

// func ValidateOnCreateCredential() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var currentReq models.RequestCreateCredential
// 		if err := ctx.ShouldBindJSON(&currentReq); err != nil {
// 			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
// 			ctx.Abort()
// 			return
// 		}
// 		ctx.Set(models.CredentialCreateContextKey, currentReq)
// 		ctx.Next()
// 	}
// }

// func ValidateOnExchangeCredential() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var currentReq models.RequestExchangeCredential
// 		if err := ctx.ShouldBindJSON(&currentReq); err != nil {
// 			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
// 			ctx.Abort()
// 			return
// 		}
// 		ctx.Set(models.CredentialExchangeContextKey, currentReq)
// 		ctx.Next()
// 	}
// }

func ValidateGetGoogleSheet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := ctx.Request.Body
		log.Printf("body: %s", body)
		var currentReq models.ActionsCommand
		if err := ctx.ShouldBindBodyWithJSON(&currentReq); err != nil {
			ctx.JSON(http.StatusBadRequest, NewInvalidRequestError(models.InvalidJSON, http.StatusBadRequest))
			ctx.Abort()
			return
		}
		ctx.Set(models.ActionGoogleKey, currentReq)
		ctx.Next()
	}
}
