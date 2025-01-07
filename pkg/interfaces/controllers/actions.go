package controllers

import (
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/repos"
	"log"
	"time"

	// "log"
	"net/http"
	// "time"

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
	// for quick response
	done := make(chan bool)
	go func() {
		// this function can return data btw
		a.actionsService.GetGoogleSheetByID(newAction.Actions)
		done <- true // Señal de que la goroutine ha terminado
	}()
	// in newaction struct, there is TestMode with true/false
	// if testMode=true data can be returned ???
	// vercel default functions 10seconds
	// currently user performs polling to retrieve data
	ctx.JSON(http.StatusAccepted, models.ResponseGetGoogleSheetByID{
		Status: http.StatusAccepted,
		Error:  "",
	})
	select {
	case <-done:
	case <-time.After(9 * time.Second):
		log.Println("WARN | Needs more time than 10 seconds")
	}
}
