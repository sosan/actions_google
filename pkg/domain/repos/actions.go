package repos

import (
	"actions_google/pkg/domain/models"
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type ActionsService interface {
	GetGoogleSheetByID(newAction *models.RequestGoogleAction) (data *[]byte)
	GetNotion(newAction *models.RequestGoogleAction) (data *[]byte)
}

type ActionsHTTPRepository interface {
	GetOAuthHTTPClient(ctx *context.Context, config *oauth2.Config, token *oauth2.Token) *http.Client
	GetDatabaseNotion(databaseID *string, secret *string) (result *models.NotionDatabaseQueryResponse, err error)
}

type ActionsRedisRepoInterface interface {
	ValidateActionGlobalUUID(field *string) (bool, error)
	SetNX(hashKey, actionID string, expiration time.Duration) (bool, error)
}

type ActionsBrokerRepository interface {
	SendAction(newAction *models.RequestGoogleAction) bool
}

type CredentialBrokerRepository interface {
	UpdateCredential(exchangeCredential *models.RequestExchangeCredential) bool
}

type ActionsNotion interface {
	ProcessNotionData(results *[]interface{}) (*[]string, *[][]string)
}
