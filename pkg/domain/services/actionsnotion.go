package services

import (
	"actions_google/pkg/common"
	"actions_google/pkg/domain/models"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

func (a *ActionsServiceImpl) getDatabaseContentFromNotion(newAction *models.RequestGoogleAction) (data *[]byte) {
	// ctx := context.Background()
	// with retries
	exchangeCredential, err := a.retriesGetCredential(newAction)
	if err != nil {
		// TODO: dead letter
		return nil
	}
	databaseID := a.getDatabaseID(newAction.Document)
	contentDB, err := a.httpRepo.GetDatabaseNotion(databaseID, &exchangeCredential.Data.Token)
	if err != nil {
		log.Printf("ERROR | %v", err)
		return nil
	}
	if len(contentDB.Results) == 0 {
		return nil
	}
	headers, arrContent := a.actionsNotion.ProcessNotionData(&contentDB.Results)
	log.Printf("%v %v", headers, arrContent)
	data = a.serializeNotionContent(headers, arrContent)
	return data
}

func (a *ActionsServiceImpl) retriesGetCredential(newAction *models.RequestGoogleAction) (*models.RequestExchangeCredential, error) {
	for i := 1; i < models.MaxAttempts; i++ {
		exchangeCredential, err := a.credentialHTTP.GetCredentialByID(&newAction.Sub, &newAction.CredentialID, 1)
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

func (a *ActionsServiceImpl) getDatabaseID(documentURI string) *string {
	splitted := strings.Split(documentURI, "/")
	idStr := strings.Split(splitted[3], "?")
	return &idStr[0]
}

func (a *ActionsServiceImpl) serializeNotionContent(headers *[]string, arrContent *[][]string) *[]byte {
	processedData := &models.ProcessedNotionData{
		Headers:     *headers,
		ContentRows: *arrContent,
	}
	jsonData, err := json.Marshal(processedData)
	if err != nil {
		log.Printf("ERROR | cannot serializeNotionContent %v %v", headers, arrContent)
		return nil
	}
	return &jsonData
}
