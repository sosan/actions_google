package services

import (
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/repos"
)

type ActionsServiceImpl struct {
	RedisRepo             repos.ActionsRedisRepoInterface
	BrokerActionsRepo     repos.ActionsBrokerRepository
	BrokerCredentialsRepo repos.CredentialBrokerRepository
	HTTPRepo              repos.ActionsHTTPRepository
	CredentialHTTP        repos.CredentialHTTPRepository
	ActionsNotion         repos.ActionsNotion
	TokenAuth             repos.TokenAuth
	SheetUtils            repos.SheetUtils
}

func NewActionsService(
	repoRedis repos.ActionsRedisRepoInterface,
	actionBroker repos.ActionsBrokerRepository,
	repoHTTP repos.ActionsHTTPRepository,
	credentialRepo repos.CredentialHTTPRepository,
	credentialBroker repos.CredentialBrokerRepository,
	notionRepo repos.ActionsNotion,
	tokenAuth repos.TokenAuth,
	sheetUtils repos.SheetUtils,

) repos.ActionsService {
	return &ActionsServiceImpl{
		RedisRepo:             repoRedis,
		BrokerActionsRepo:     actionBroker,
		BrokerCredentialsRepo: credentialBroker,
		HTTPRepo:              repoHTTP,
		CredentialHTTP:        credentialRepo,
		ActionsNotion:         notionRepo,
		TokenAuth:             tokenAuth,
		SheetUtils:            sheetUtils,
	}
}

func (a *ActionsServiceImpl) GetGoogleSheetByID(newAction *models.RequestGoogleAction) (data *[]byte) {
	if newAction == nil {
		return nil
	}
	// retries???
	switch newAction.Operation {
	case "getallcontent":
		data = a.GetAllContentFromGoogleSheets(newAction)
	default:
		return nil
	}
	if data == nil || string(*data) == "" {
		return nil
	}
	newAction.Data = string(*data)
	a.BrokerActionsRepo.SendAction(newAction)
	return data
}

func (a *ActionsServiceImpl) GetNotion(newAction *models.RequestGoogleAction) (data *[]byte) {
	if newAction == nil {
		return nil
	}
	// retries???
	switch newAction.Operation {
	case "getallcontent":
		data = a.getDatabaseContentFromNotion(newAction)
	default:
		return nil
	}
	if data == nil || string(*data) == "" {
		return nil
	}
	newAction.Data = string(*data)
	a.BrokerActionsRepo.SendAction(newAction)
	return data
}
