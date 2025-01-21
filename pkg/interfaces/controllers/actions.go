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
		// this function return byte(data)
		a.actionsService.GetGoogleSheetByID(newAction.Actions)
		done <- true // goroutine ended
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
	case <-time.After(models.MaxSecondsGoRoutine):
		log.Println("WARN | Needs more time than 10 seconds")
	}
}

// TODO: right now not consult type of action, only implemented read
func (a *ActionsController) GetNotion(ctx *gin.Context) {
	newAction := ctx.MustGet(models.ActionNotionKey).(models.ActionsCommand)
	// for quick response
	done := make(chan bool)
	go func() {
		// this function returns byte data, not necessary
		a.actionsService.GetNotion(newAction.Actions)
		done <- true // goroutine ended
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
	case <-time.After(models.MaxSecondsGoRoutine):
		log.Println("WARN | Needs more time than 10 seconds")
	}
}
