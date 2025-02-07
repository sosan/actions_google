package dimodel

import (
	"actions_google/pkg/domain/repos"
	"actions_google/pkg/interfaces/controllers"
)

type Dependencies struct {
	AuthService       *repos.AuthService
	AuthController    *controllers.AuthController
	ActionsController *controllers.ActionsController
}
