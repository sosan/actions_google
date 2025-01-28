package repos

import (
	"actions_google/pkg/domain/models"
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/sheets/v4"
)

type ActionsService interface {
	GetGoogleSheetByID(newAction *models.RequestGoogleAction) (data *[]byte)
	GetNotion(newAction *models.RequestGoogleAction) (data *[]byte)
	GetAllContentFromGoogleSheets(newAction *models.RequestGoogleAction) (data *[]byte)
}

type TokenAuth interface {
	GenerateTokenOAuth(*context.Context, *oauth2.Config, *models.RequestExchangeCredential) *oauth2.Token
	GetConfigOAuth(models.DataCredential) *oauth2.Config
}

type SheetUtils interface {
	GetAllContentFromGoogleSheets(document *string, client *http.Client, actionID *string) (*sheets.ValueRange, error)
	GetSpreadsheetID(documentURI *string) *string
	CreateSheetsService(ctx context.Context, client *http.Client) (*sheets.Service, error)
	GetValuesFromSheet(sheets *sheets.Spreadsheet, sheetsService *sheets.Service, spreadsheetID *string) (*sheets.ValueRange, error)
	GetSpreadsheet(ctx context.Context, srv *sheets.Service, spreadsheetID string) (*sheets.Spreadsheet, error)
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
	UpdateCredential(exchangeCredential *models.RequestExchangeCredential, token *oauth2.Token) bool
	UpdateCredentialFromGoogle(exchangeCredential *models.RequestExchangeCredential, token *oauth2.Token) *models.RequestExchangeCredential
}

type TransformNotion interface {
	ProcessNotionData(results *[]interface{}) (*[]string, *[][]string)
}
