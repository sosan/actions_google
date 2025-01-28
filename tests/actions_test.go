package tests

import (
	"actions_google/mocks"
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/services"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type MockStructured struct {
	CredentialHTTP        *mocks.CredentialHTTPRepository
	BrokerCredentialsRepo *mocks.CredentialBrokerRepository
	HTTPRepo              *mocks.ActionsHTTPRepository
	TokenAuth             *mocks.TokenAuth
	SheetUtils            *mocks.SheetUtils
}

func createNewMocks() *MockStructured {
	ctrl := &MockStructured{
		CredentialHTTP:        new(mocks.CredentialHTTPRepository),
		BrokerCredentialsRepo: new(mocks.CredentialBrokerRepository),
		HTTPRepo:              new(mocks.ActionsHTTPRepository),
		TokenAuth:             new(mocks.TokenAuth),
		SheetUtils:            new(mocks.SheetUtils),
	}
	return ctrl
}

func TestGetAllContentFromGoogleSheets(t *testing.T) {

	testCredential := &models.RequestExchangeCredential{
		Data: models.DataCredential{
			ClientID:     "mock-client",
			ClientSecret: "mock-secret",
			TokenRefresh: "mock-refresh-token",
		},
	}

	testOauthConfig := &oauth2.Config{
		RedirectURL:  "data.RedirectURL",
		ClientID:     "data.ClientID",
		ClientSecret: "data.ClientSecret",
		Scopes:       []string{"data.Scopes"},
		Endpoint:     google.Endpoint,
	}

	testToken := &oauth2.Token{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	t.Run("Successfully retrieved data from Sheets", func(t *testing.T) {
		ctrl := createNewMocks()

		a := &services.ActionsServiceImpl{
			CredentialHTTP:        ctrl.CredentialHTTP,
			BrokerCredentialsRepo: ctrl.BrokerCredentialsRepo,
			HTTPRepo:              ctrl.HTTPRepo,
			TokenAuth:             ctrl.TokenAuth,
			SheetUtils:            ctrl.SheetUtils,
		}

		expectedDoc := "https://www.ejemplo.com/1/2/3/xmock-spreadsheet-id"
		expectedActionID := "accion-123"
		expectedClient := &http.Client{}

		ctrl.TokenAuth.On("GenerateTokenOAuth", mock.Anything, mock.Anything, mock.Anything).Return(testToken)
		ctrl.TokenAuth.On("GetConfigOAuth", mock.Anything).Return(testOauthConfig)
		ctrl.SheetUtils.On("GetSpreadsheetID", mock.Anything).Return("mock-spreadsheet-id")
		ctrl.HTTPRepo.On("GetOAuthHTTPClient", mock.Anything, mock.Anything, mock.Anything).Return(expectedClient)
		ctrl.CredentialHTTP.On(
			"GetCredentialByID",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(testCredential, nil)
		ctrl.BrokerCredentialsRepo.On(
			"UpdateCredential",
			mock.Anything,
			testToken,
		).Return(true)

		ctrl.SheetUtils.On(
			"GetAllContentFromGoogleSheets",
			mock.MatchedBy(func(doc *string) bool {
				return *doc == expectedDoc
			}),
			expectedClient,
			mock.MatchedBy(func(id *string) bool {
				return *id == expectedActionID
			}),
		).Return(&sheets.ValueRange{Values: [][]interface{}{{"A1", "B1"}}}, nil)

		newAction := &models.RequestGoogleAction{
			Document:   expectedDoc,
			ActionID:   expectedActionID,
			Sub:        "test-sub",
			WorkflowID: "test-workflow",
			NodeID:     "test-node",
		}

		result := a.GetAllContentFromGoogleSheets(newAction)

		assert.NotNil(t, result)
		ctrl.SheetUtils.AssertCalled(t, "GetAllContentFromGoogleSheets", mock.Anything, expectedClient, mock.Anything)
	})

	t.Run("Error fetching credentials", func(t *testing.T) {
		ctrl := createNewMocks()

		a := &services.ActionsServiceImpl{
			CredentialHTTP:        ctrl.CredentialHTTP,
			BrokerCredentialsRepo: ctrl.BrokerCredentialsRepo,
			HTTPRepo:              ctrl.HTTPRepo,
			TokenAuth:             ctrl.TokenAuth,
			SheetUtils:            ctrl.SheetUtils,
		}

		ctrl.CredentialHTTP.On(
			"GetCredentialByID",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, assert.AnError)

		newAction := &models.RequestGoogleAction{
			Document:   "https://url-valida",
			ActionID:   "accion-123",
			Sub:        "test-sub",
			WorkflowID: "test-workflow",
			NodeID:     "test-node",
		}

		result := a.GetAllContentFromGoogleSheets(newAction)
		assert.Nil(t, result)
		ctrl.CredentialHTTP.AssertExpectations(t)
	})

	t.Run("Token generation error", func(t *testing.T) {
		ctrl := createNewMocks()

		a := &services.ActionsServiceImpl{
			CredentialHTTP:        ctrl.CredentialHTTP,
			BrokerCredentialsRepo: ctrl.BrokerCredentialsRepo,
			HTTPRepo:              ctrl.HTTPRepo,
			TokenAuth:             ctrl.TokenAuth,
			SheetUtils:            ctrl.SheetUtils,
		}

		ctrl.CredentialHTTP.On("GetCredentialByID", mock.Anything, mock.Anything, mock.Anything).Return(testCredential, nil)
		ctrl.TokenAuth.On("GetConfigOAuth", mock.Anything).Return(testOauthConfig)
		ctrl.TokenAuth.On("GenerateTokenOAuth", mock.Anything, mock.Anything, mock.Anything).Return(nil) // Token nil

		newAction := &models.RequestGoogleAction{
			Document:   "https://url-valida",
			ActionID:   "accion-123",
			Sub:        "test-sub",
			WorkflowID: "test-workflow",
			NodeID:     "test-node",
		}

		result := a.GetAllContentFromGoogleSheets(newAction)
		assert.Nil(t, result)
		ctrl.TokenAuth.AssertExpectations(t)
	})

	t.Run("Error creating HTTP client", func(t *testing.T) {
		ctrl := createNewMocks()

		a := &services.ActionsServiceImpl{
			CredentialHTTP:        ctrl.CredentialHTTP,
			BrokerCredentialsRepo: ctrl.BrokerCredentialsRepo,
			HTTPRepo:              ctrl.HTTPRepo,
			TokenAuth:             ctrl.TokenAuth,
			SheetUtils:            ctrl.SheetUtils,
		}

		ctrl.CredentialHTTP.On("GetCredentialByID", mock.Anything, mock.Anything, mock.Anything).Return(testCredential, nil)
		ctrl.TokenAuth.On("GenerateTokenOAuth", mock.Anything, mock.Anything, mock.Anything).Return(testToken)
		ctrl.TokenAuth.On("GetConfigOAuth", mock.Anything).Return(testOauthConfig)
		ctrl.HTTPRepo.On("GetOAuthHTTPClient", mock.Anything, mock.Anything, mock.Anything).Return(nil) // Client nil

		newAction := &models.RequestGoogleAction{
			Document:   "https://url-valida",
			ActionID:   "accion-123",
			Sub:        "test-sub",
			WorkflowID: "test-workflow",
			NodeID:     "test-node",
		}

		result := a.GetAllContentFromGoogleSheets(newAction)
		assert.Nil(t, result)
		ctrl.HTTPRepo.AssertExpectations(t)
	})

	t.Run("Error updating credentials", func(t *testing.T) {
		ctrl := createNewMocks()

		a := &services.ActionsServiceImpl{
			CredentialHTTP:        ctrl.CredentialHTTP,
			BrokerCredentialsRepo: ctrl.BrokerCredentialsRepo,
			HTTPRepo:              ctrl.HTTPRepo,
			TokenAuth:             ctrl.TokenAuth,
			SheetUtils:            ctrl.SheetUtils,
		}

		ctrl.CredentialHTTP.On("GetCredentialByID", mock.Anything, mock.Anything, mock.Anything).Return(testCredential, nil)
		ctrl.TokenAuth.On("GenerateTokenOAuth", mock.Anything, mock.Anything, mock.Anything).Return(testToken)
		ctrl.TokenAuth.On("GetConfigOAuth", mock.Anything).Return(testOauthConfig)
		ctrl.HTTPRepo.On("GetOAuthHTTPClient", mock.Anything, mock.Anything, mock.Anything).Return(&http.Client{})
		ctrl.SheetUtils.On("GetSpreadsheetID", mock.Anything).Return("spreadsheet-valido")
		ctrl.SheetUtils.On(
			"GetAllContentFromGoogleSheets",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&sheets.ValueRange{Values: [][]interface{}{{"A1"}}}, nil)
		ctrl.BrokerCredentialsRepo.On("UpdateCredential", mock.Anything, mock.Anything).Return(false) // Failed to update

		newAction := &models.RequestGoogleAction{
			Document:   "https://url-valida",
			ActionID:   "accion-123",
			Sub:        "test-sub",
			WorkflowID: "test-workflow",
			NodeID:     "test-node",
		}

		result := a.GetAllContentFromGoogleSheets(newAction)
		assert.NotNil(t, result)
		ctrl.BrokerCredentialsRepo.AssertExpectations(t)
	})

	t.Run("SheetUtils error", func(t *testing.T) {
		ctrl := createNewMocks()

		a := &services.ActionsServiceImpl{
			CredentialHTTP:        ctrl.CredentialHTTP,
			BrokerCredentialsRepo: ctrl.BrokerCredentialsRepo,
			HTTPRepo:              ctrl.HTTPRepo,
			TokenAuth:             ctrl.TokenAuth,
			SheetUtils:            ctrl.SheetUtils,
		}

		expectedDoc := "https://url-valida"
		expectedActionID := "accion-123"

		ctrl.SheetUtils.On(
			"GetAllContentFromGoogleSheets",
			mock.MatchedBy(func(doc *string) bool {
				return *doc == expectedDoc
			}),
			mock.Anything, // HTTP Client
			mock.MatchedBy(func(id *string) bool {
				return *id == expectedActionID
			}),
		).Return(nil, assert.AnError)

		ctrl.CredentialHTTP.On("GetCredentialByID", mock.Anything, mock.Anything, mock.Anything).Return(testCredential, nil)
		ctrl.TokenAuth.On("GenerateTokenOAuth", mock.Anything, mock.Anything, mock.Anything).Return(testToken)
		ctrl.TokenAuth.On("GetConfigOAuth", mock.Anything).Return(testOauthConfig)
		ctrl.HTTPRepo.On("GetOAuthHTTPClient", mock.Anything, mock.Anything, mock.Anything).Return(&http.Client{})

		newAction := &models.RequestGoogleAction{
			Document:   expectedDoc,
			ActionID:   expectedActionID,
			Sub:        "test-sub",
			WorkflowID: "test-workflow",
			NodeID:     "test-node",
		}

		result := a.GetAllContentFromGoogleSheets(newAction)
		assert.Nil(t, result)
		ctrl.SheetUtils.AssertCalled(t, "GetAllContentFromGoogleSheets", mock.Anything, mock.Anything, mock.Anything)
	})
}
