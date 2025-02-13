package services

import (
	"actions_google/pkg/domain/models"
	"context"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func (a *ActionsServiceImpl) GetAllContentFromGoogleSheets(newAction *models.RequestGoogleAction) (data *[]byte) {
	ctx := context.Background()
	exchangeCredential, err := a.RetriesGetCredential(newAction)
	if err != nil {
		log.Printf("ERROR | Cannot fetching credential by ID: %v", err)
		// TODO: dead letter
		return nil
	}

	// exchangeCredential, err := a.CredentialHTTP.GetCredentialByID(&newAction.Sub, &newAction.CredentialID, 1)
	// if err != nil {
	// 	return nil
	// }
	config := a.TokenAuth.GetConfigOAuth(exchangeCredential.Data)
	// this new token needs to be updated to DB
	token := a.TokenAuth.GenerateTokenOAuth(&ctx, config, exchangeCredential)
	if token == nil {
		// TODO: deadletter
		log.Printf("ERROR | Failed to generate OAuth token for user %s workflowid %s nodeid %s actionid %s", newAction.Sub, newAction.WorkflowID, newAction.NodeID, newAction.ActionID)
		return nil
	}
	httpClient := a.getClient(&ctx, config, token)
	if httpClient == nil {
		log.Printf("ERROR | Failed to create HTTP client for user %s workflowid %s nodeid %s actionid %s", newAction.Sub, newAction.WorkflowID, newAction.NodeID, newAction.ActionID)
		return nil
	}

	values, err := a.SheetUtils.GetAllContentFromGoogleSheets(&newAction.Document, httpClient, &newAction.ActionID)
	if err != nil {
		return nil
	}

	if values == nil {
		log.Printf("ERROR | No values found")
		return nil
	}

	// Save new token and refrestoken to DB
	// this operation CAN FAIL to save to DB NOT implemented retries and deadletters
	updated := a.BrokerCredentialsRepo.UpdateCredential(exchangeCredential, token)
	if !updated {
		log.Printf("WARN | Failed to update credentials in the database for CredentialID: %s", exchangeCredential.ID)
		// TODO: retries
		// TODO: dead letter
	}
	// ---
	// log.Printf("%v", values)
	str, err := values.MarshalJSON()
	if err != nil {
		log.Printf("ERROR | marshalling values to JSON: %v for actionid: %s", err, newAction.ActionID)
		return nil
	}

	return &str
}

// TODO: repo httpclient
func (a *ActionsServiceImpl) getClient(ctx *context.Context, config *oauth2.Config, token *oauth2.Token) *http.Client {
	client := a.HTTPRepo.GetOAuthHTTPClient(ctx, config, token)
	return client
}
