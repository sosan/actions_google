package controllers

import (
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/repos"
	"log"
	// "net/http"

	"github.com/gin-gonic/gin"
)

type ActionsController struct {
	actionsService repos.ActionsService
}

func NewActionsController(newActionsService repos.ActionsService) *ActionsController {
	return &ActionsController{actionsService: newActionsService}
}

func (a *ActionsController) Ping(ctx *gin.Context) {
	ob := gin.H{
		"test": "test",
	}
	ctx.JSON(200, ob)
}

func (a *ActionsController) GetGoogleSheetByID(ctx *gin.Context) {
	newAction := ctx.MustGet(models.ActionGoogleKey).(models.ActionsCommand)
	log.Printf("%v", newAction)
	// created, exist, actionsData := a.actionsService.GetGoogleSheetByID(newAction.Actions)
	// if !created && !exist {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error":  "not generated",
	// 		"status": http.StatusInternalServerError,
	// 	})
	// 	return
	// }

	// if exist {
	// 	ctx.JSON(http.StatusAlreadyReported, gin.H{
	// 		"error":  "asdasdadasd",
	// 		"status": http.StatusAlreadyReported,
	// 	})
	// 	return
	// }

	// ctx.JSON(http.StatusCreated, models.ResponseGetGoogleSheetByID{
	// 	Status: http.StatusOK,
	// 	Error:  "",
	// 	Action: *actionsData,
	// })
}
