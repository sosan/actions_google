package tests

import (
	// "actions_google/mocks"
	models "actions_google/pkg/domain/models"
	"actions_google/pkg/domain/services"
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestGetDatabaseContentFromNotion_Success(t *testing.T) {
	databaseID := uuid.New().String()
	token := "test-token"
	documentURI := "https://www.notion.so/" + databaseID + "?v=test"
	newAction := &models.RequestGoogleAction{
		Document:     documentURI,
		CredentialID: "cred-123",
		Sub:          "user-123",
	}

	exchangeCredential := &models.RequestExchangeCredential{
		Data: models.DataCredential{Token: token},
	}
	contentDB := &models.NotionDatabaseQueryResponse{
		Results: []interface{}{map[string]interface{}{"prop": "value"}},
	}
	headers := []string{"Header1"}
	content := [][]string{{"Value1"}}

	t.Run("Successfully retrieved data from Sheets", func(t *testing.T) {
		ctrl := createNewMocks()

		a := &services.ActionsServiceImpl{
			CredentialHTTP:        ctrl.CredentialHTTP,
			BrokerCredentialsRepo: ctrl.BrokerCredentialsRepo,
			HTTPRepo:              ctrl.HTTPRepo,
			TokenAuth:             ctrl.TokenAuth,
			SheetUtils:            ctrl.SheetUtils,
			ActionsNotion:         ctrl.ActionsNotion,
		}

		ctrl.ActionsNotion.On("RetriesGetCredential", newAction).Return(exchangeCredential, nil)
		ctrl.BrokerCredentialsRepo.On("GetCredential", "cred-123", "user-123").Return(exchangeCredential, nil)
		ctrl.HTTPRepo.On("GetDatabaseNotion", &databaseID, &token).Return(contentDB, nil)
		ctrl.CredentialHTTP.On("GetCredentialByID", &newAction.Sub, &newAction.CredentialID, uint64(1)).Return(exchangeCredential, nil)
		ctrl.ActionsNotion.On("ProcessNotionData", &contentDB.Results).Return(&headers, &content)

		expectedData, _ := json.Marshal(&models.ProcessedNotionData{
			Headers:     headers,
			ContentRows: content,
		})

		gotData := a.GetDatabaseContentFromNotion(newAction)
		if gotData == nil || !bytes.Equal(*gotData, expectedData) {
			t.Errorf("Not same, got: %v exepected: %v", gotData, expectedData)
		}
	})

	t.Run("Empty results with retries (too much time)", func(t *testing.T) {
		ctrl := createNewMocks()

		a := &services.ActionsServiceImpl{
			CredentialHTTP:        ctrl.CredentialHTTP,
			BrokerCredentialsRepo: ctrl.BrokerCredentialsRepo,
			HTTPRepo:              ctrl.HTTPRepo,
			TokenAuth:             ctrl.TokenAuth,
			SheetUtils:            ctrl.SheetUtils,
			ActionsNotion:         ctrl.ActionsNotion,
		}

		ctrl.ActionsNotion.On("RetriesGetCredential", newAction).Return(nil, nil)
		ctrl.BrokerCredentialsRepo.On("GetCredential", "cred-123", "user-123").Return(nil, nil)
		ctrl.HTTPRepo.On("GetDatabaseNotion", &databaseID, &token).Return(contentDB, nil)
		ctrl.CredentialHTTP.On("GetCredentialByID", &newAction.Sub, &newAction.CredentialID, uint64(1)).Return(nil, nil)
		ctrl.ActionsNotion.On("ProcessNotionData", &contentDB.Results).Return(&headers, &content)

		var expectedData []byte = nil

		gotData := a.GetDatabaseContentFromNotion(newAction)
		if gotData != nil {
			t.Errorf("Not same, got: %v exepected: %v", gotData, expectedData)
		}
	})

	t.Run("HTTP Error", func(t *testing.T) {
		ctrl := createNewMocks()

		a := &services.ActionsServiceImpl{
			CredentialHTTP:        ctrl.CredentialHTTP,
			BrokerCredentialsRepo: ctrl.BrokerCredentialsRepo,
			HTTPRepo:              ctrl.HTTPRepo,
			TokenAuth:             ctrl.TokenAuth,
			SheetUtils:            ctrl.SheetUtils,
			ActionsNotion:         ctrl.ActionsNotion,
		}

		ctrl.ActionsNotion.On("RetriesGetCredential", newAction).Return(exchangeCredential, nil)
		ctrl.BrokerCredentialsRepo.On("GetCredential", "cred-123", "user-123").Return(exchangeCredential, nil)
		ctrl.HTTPRepo.On("GetDatabaseNotion", &databaseID, &token).Return(nil, errors.New("HTTP error"))
		ctrl.CredentialHTTP.On("GetCredentialByID", &newAction.Sub, &newAction.CredentialID, uint64(1)).Return(exchangeCredential, nil)
		ctrl.ActionsNotion.On("ProcessNotionData", &contentDB.Results).Return(&headers, &content)

		expectedData, _ := json.Marshal(&models.ProcessedNotionData{
			Headers:     headers,
			ContentRows: content,
		})

		gotData := a.GetDatabaseContentFromNotion(newAction)
		if gotData != nil {
			t.Errorf("Not same, got: %v exepected: %v", gotData, expectedData)
		}
	})
}
