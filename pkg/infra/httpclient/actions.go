package httpclient

import (
	"actions_google/pkg/config"
	"actions_google/pkg/domain/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type ActionsHTTPRepository struct {
	databaseHTTPURL string
	token           string
	client          HTTPClient
}

func NewActionsClientHTTP(client HTTPClient, clickhouseConfig config.ClickhouseConfig) *ActionsHTTPRepository {
	return &ActionsHTTPRepository{
		client:          client,
		databaseHTTPURL: clickhouseConfig.GetClickhouseURI(),
		token:           clickhouseConfig.GetClickhouseToken(),
	}
}

func (a *ActionsHTTPRepository) GetOAuthHTTPClient(ctx *context.Context, config *oauth2.Config, token *oauth2.Token) *http.Client {
	client := config.Client(*ctx, token)
	return client
}

func (a *ActionsHTTPRepository) GetDatabaseNotion(databaseID *string, secret *string) (result *models.NotionDatabaseQueryResponse, err error) {
	// without context
	apiURL := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", *databaseID)
	validatedURI, err := a.validateURL(apiURL)
	if err != nil {
		log.Printf("ERROR | cannot parse url %s", apiURL)
		return nil, err
	}
	req, err := http.NewRequest("POST", validatedURI, nil)
	if err != nil {
		return nil, err
	}

	a.setHeadersNotion(req, secret)

	resp, err := a.client.Do(req)
	if err != nil {
		log.Printf("ERROR | actionsgoogle cannot send request %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ERROR | failed to retrieve database response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
		return nil, fmt.Errorf("ERROR | cannot decode token: %v", err)
	}
	return result, nil
}

func (a *ActionsHTTPRepository) validateURL(rawURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL")
	}

	return parsedURL.String(), nil
}

func (a *ActionsHTTPRepository) setHeadersNotion(req *http.Request, secret *string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *secret))
	// req.Header.Set("Notion-Version", "2023-09-22") // latest
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")
}
