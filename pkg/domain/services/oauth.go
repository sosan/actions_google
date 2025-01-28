package services

import (
	"actions_google/pkg/domain/models"
	"context"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type TokenAuthImpl struct {
}

func NewTokenAuthImpl() *TokenAuthImpl {
	return &TokenAuthImpl{}
}

// function to generate a new refresh token
// once the refresh token is generated,
// DB needed to be updated
func (t *TokenAuthImpl) GenerateTokenOAuth(
	ctx *context.Context,
	config *oauth2.Config,
	cred *models.RequestExchangeCredential,
) *oauth2.Token {
	token := &oauth2.Token{RefreshToken: cred.Data.TokenRefresh}
	tokenSource := config.TokenSource(*ctx, token)
	// if OAuth consent screen is in dev mode, those
	// refresh tokens it will expire in 7 days
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Printf("ERROR | cannot renovate token: config %v %v", config, cred)
		return nil
	}
	return newToken
}

func (t *TokenAuthImpl) GetConfigOAuth(data models.DataCredential) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  data.RedirectURL,
		ClientID:     data.ClientID,
		ClientSecret: data.ClientSecret,
		Scopes:       data.Scopes,
		Endpoint:     google.Endpoint,
	}
}
