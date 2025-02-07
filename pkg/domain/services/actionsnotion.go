package services

import (
	"actions_google/pkg/domain/models"
	"encoding/json"

	"log"
	"strings"

	"github.com/google/uuid"
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

// GetDatabaseID returns the ID of the database from the document URI
// Notion Pattern URI https://www.notion.so/{workspace}/{databaseID}?v={version}
// Notion Pattern URI https://www.notion.so/{databaseID}?v={version}
// in validation we checked if documentURI is empty
// in validation we checked if documentURI contains NotionHost
func (a *ActionsServiceImpl) GetDatabaseID(documentURI string) *string {
	if strings.TrimSpace(documentURI) == "" {
		return nil
	}

	splitted := strings.SplitN(documentURI, models.NotionHost, 2)
	if len(splitted) < 2 {
		return nil
	}

	pathAndQuery := strings.SplitN(splitted[1], "?", 2)[0]
	pathSegments := strings.Split(strings.Trim(pathAndQuery, "/"), "/")

	if len(pathSegments) == 0 {
		return nil
	}

	databaseID := pathSegments[len(pathSegments)-1]
	if _, err := uuid.Parse(databaseID); err != nil {
		log.Printf("ERROR | Cannot parse ID: %v", err)
		return nil
	}

	return &databaseID
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
