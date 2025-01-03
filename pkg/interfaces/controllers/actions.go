package controllers

import (
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/repos"
	"net/http"

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
	data := a.actionsService.GetGoogleSheetByID(newAction.Actions)
	// data not used
	// there is an bool option devtest can be used to response directly
	if string(*data) == "" {
		ctx.JSON(http.StatusInternalServerError, models.ResponseGetGoogleSheetByID{
			Error:  "not generated",
			Status: http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseGetGoogleSheetByID{
		Status: http.StatusOK,
		Error:  "",
		// Data:   newAction.Actions.ActionID, //not necesary
	})
}
