package services

import (
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/repos"
)

type ActionsServiceImpl struct {
	redisRepo             repos.ActionsRedisRepoInterface
	brokerActionsRepo     repos.ActionsBrokerRepository
	brokerCredentialsRepo repos.CredentialBrokerRepository
	httpRepo              repos.ActionsHTTPRepository
	credentialHTTP        repos.CredentialHTTPRepository
}

func NewActionsService(repoRedis repos.ActionsRedisRepoInterface, actionBroker repos.ActionsBrokerRepository, repoHTTP repos.ActionsHTTPRepository, credentialRepo repos.CredentialHTTPRepository, credentialBroker repos.CredentialBrokerRepository) repos.ActionsService {
	return &ActionsServiceImpl{
		redisRepo:             repoRedis,
		brokerActionsRepo:     actionBroker,
		brokerCredentialsRepo: credentialBroker,
		httpRepo:              repoHTTP,
		credentialHTTP:        credentialRepo,
	}
}

func (a *ActionsServiceImpl) GetGoogleSheetByID(newAction *models.RequestGoogleAction) (data *[]byte) {
	if newAction == nil {
		return nil
	}
	// retries???
	switch newAction.Operation {
	case "getallcontent":
		data = a.getAllContentFromGoogleSheets(newAction)
	default:
		return nil
	}
	if data == nil || string(*data) == "" {
		return nil
	}
	newAction.Data = string(*data)
	a.brokerActionsRepo.SendAction(newAction)
	return data
}
