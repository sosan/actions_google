package services

import (
	"actions_google/pkg/common"
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/repos"
	"fmt"
	"log"
	"time"
)

type ActionsServiceImpl struct {
	RedisRepo             repos.ActionsRedisRepoInterface
	BrokerActionsRepo     repos.ActionsBrokerRepository
	BrokerCredentialsRepo repos.CredentialBrokerRepository
	HTTPRepo              repos.ActionsHTTPRepository
	CredentialHTTP        repos.CredentialHTTPRepository
	ActionsNotion         repos.TransformNotion
	TokenAuth             repos.TokenAuth
	SheetUtils            repos.SheetUtils
}

func NewActionsService(
	repoRedis repos.ActionsRedisRepoInterface,
	actionBroker repos.ActionsBrokerRepository,
	repoHTTP repos.ActionsHTTPRepository,
	credentialRepo repos.CredentialHTTPRepository,
	credentialBroker repos.CredentialBrokerRepository,
	notionRepo repos.TransformNotion,
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
		data = a.GetDatabaseContentFromNotion(newAction)
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

func (a *ActionsServiceImpl) RetriesGetCredential(newAction *models.RequestGoogleAction) (*models.RequestExchangeCredential, error) {
	for i := 1; i < models.MaxAttempts; i++ {
		exchangeCredential, err := a.CredentialHTTP.GetCredentialByID(&newAction.Sub, &newAction.CredentialID, 1)
		if err != nil {
			log.Printf("ERROR | Cannot fetching credential by ID: %v", err)
			return nil, err
		}
		if exchangeCredential != nil {
			return exchangeCredential, err
		}
		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("WARNING | Failed to create action %s for user %s , attempt %d:. Retrying in %v", newAction.ActionID, newAction.Sub, i, waitTime)
		time.Sleep(waitTime)
	}
	return nil, fmt.Errorf("cannot fetching credential by sub %s credentialid %s", newAction.Sub, newAction.CredentialID)
}
