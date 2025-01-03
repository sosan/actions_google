package services

import (
	"actions_google/pkg/domain/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (a *ActionsServiceImpl) getAllContentFromGoogleSheets(newAction *models.RequestGoogleAction) (data *[]byte) {
	ctx := context.Background()
	exchangeCredential, err := a.credentialHTTP.GetCredentialByID(&newAction.Sub, &newAction.CredentialID, 1)
	if err != nil {
		return nil
	}
	config := a.getConfigOAuth(exchangeCredential.Data)
	// this new token needs to be updated to DB
	token := a.generateTokenOAuth(&ctx, config, exchangeCredential)
	if token == nil {
		return nil
	}
	httpClient := a.getClient(&ctx, config, token)
	if httpClient == nil {
		return nil
	}
	sheetsService, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil
	}
	spreadsheetID := a.getSpreedSheetID(&newAction.Document)
	response, err := sheetsService.Spreadsheets.Get(*spreadsheetID).Context(ctx).Do()
	if response == nil || err != nil {
		return nil
	}
	values, err := a.getValuesFromSheet(response, sheetsService, spreadsheetID)
	if values == nil || err != nil {
		return nil
	}

	// TODO: it can be moved to another position
	// not necessary returned value btw for more redeability
	exchangeCredential = a.updateCredentialFromGoogle(exchangeCredential, token)
	// Save new token and refrestoken to DB
	// this operation CAN FAIL to save to DB NOT implemented retries and deadletters
	updated := a.brokerCredentialsRepo.UpdateCredential(exchangeCredential)
	if !updated {
		// TODO: retries
		// TODO: dead letter
		log.Printf("updated %v", updated )
	}
	// ---
	// log.Printf("%v", values)
	str, err := values.MarshalJSON()
	if err != nil {
		return nil
	}

	return &str
}

func (a *ActionsServiceImpl) getConfigOAuth(data models.DataCredential) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  data.RedirectURL,
		ClientID:     data.ClientID,
		ClientSecret: data.ClientSecret,
		Scopes:       data.Scopes,
		Endpoint:     google.Endpoint,
	}
}

// function to generate a new refresh token
// once the refresh token is generated,
// DB needed to be updated
func (a *ActionsServiceImpl) generateTokenOAuth(ctx *context.Context, config *oauth2.Config, credential *models.RequestExchangeCredential) *oauth2.Token {
	token := &oauth2.Token{
		RefreshToken: credential.Data.TokenRefresh,
	}
	tokenSource := config.TokenSource(*ctx, token)
	// if OAuth consent screen is in dev mode, those
	// refresh tokens it will expire in 7 days
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Printf("ERROR | cannot renovate token: config %v %v", config, credential)
		return nil
	}
	return newToken
}

// TODO: repo httpclient
func (a *ActionsServiceImpl) getClient(ctx *context.Context, config *oauth2.Config, token *oauth2.Token) *http.Client {
	client := a.httpRepo.GetOAuthHTTPClient(ctx, config, token)
	return client
}

// https://docs.google.com/spreadsheets/d/1o8Znm0MXXX0MDsMMMNNfnB7Q7hs2T08MMYnpbQANchs/edit?gid=0#gid=0
// https://docs.google.com/spreadsheets/d/1BxiMVs0XXX5nFMMMMBdBZjgmUUqptlbs74OgvE2upms/edit
func (a *ActionsServiceImpl) getSpreedSheetID(documentURI *string) *string {
	id := strings.Split(*documentURI, "/")[5]
	return &id
}

func (a *ActionsServiceImpl) getValuesFromSheet(sheets *sheets.Spreadsheet, sheetsService *sheets.Service, spreadsheetID *string) (*sheets.ValueRange, error) {
	for _, sheet := range sheets.Sheets {
		// properties := sheet.Properties
		// gridProperties := properties.GridProperties
		// dinamic range
		// readRange := fmt.Sprintf("%s", sheetName)
		// readRange := fmt.Sprintf("%s!A1:%s%d", sheetName, getColumnName(gridProperties.ColumnCount), gridProperties.RowCount)
		sheetName := sheet.Properties.Title
		readRange := sheetName

		values, err := sheetsService.Spreadsheets.Values.Get(*spreadsheetID, readRange).Do()
		return values, err
	}
	return nil, fmt.Errorf("ERROR | not sheets")
}

func (a *ActionsServiceImpl) updateCredentialFromGoogle(exchangeCredential *models.RequestExchangeCredential, token *oauth2.Token) *models.RequestExchangeCredential {
	exchangeCredential.Data.Token = token.AccessToken
	exchangeCredential.Data.TokenRefresh = token.RefreshToken
	exchangeCredential.UpdatedAt.Time = time.Now().UTC()

	// token.expiry already set to 0
	exchangeCredential.ExpiresAt.Time = token.Expiry.UTC().Add(-models.TimeDriftForExpire * time.Second)
	return exchangeCredential
}
