package httpclient

import (
	"actions_google/pkg/config"
	"context"
	"net/http"

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
