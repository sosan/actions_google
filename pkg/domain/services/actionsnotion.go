package services

import (
	
	"actions_google/pkg/domain/models"
	"encoding/json"
	
	"log"
	"strings"
	
)

func (a *ActionsServiceImpl) GetDatabaseContentFromNotion(newAction *models.RequestGoogleAction) (data *[]byte) {
	// ctx := context.Background()
	// with retries
	exchangeCredential, err := a.RetriesGetCredential(newAction)
	if err != nil {
		// TODO: dead letter
		return nil
	}
	databaseID := a.GetDatabaseID(newAction.Document)
	contentDB, err := a.HTTPRepo.GetDatabaseNotion(databaseID, &exchangeCredential.Data.Token)
	if err != nil {
		log.Printf("ERROR | %v", err)
		return nil
	}
	if len(contentDB.Results) == 0 {
		return nil
	}
	headers, arrContent := a.ActionsNotion.ProcessNotionData(&contentDB.Results)
	log.Printf("%v %v", headers, arrContent)
	data = a.SerializeNotionContent(headers, arrContent)
	return data
}



func (a *ActionsServiceImpl) GetDatabaseID(documentURI string) *string {
	splitted := strings.Split(documentURI, "/")
	idStr := strings.Split(splitted[3], "?")
	return &idStr[0]
}

func (a *ActionsServiceImpl) SerializeNotionContent(headers *[]string, arrContent *[][]string) *[]byte {
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
