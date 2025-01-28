package brokerclient

import (
	"encoding/json"
	// "fmt"
	"actions_google/pkg/common"
	"actions_google/pkg/domain/models"
	"log"
	"time"

	"golang.org/x/oauth2"
)

type CredentialCommand struct {
	Type       string                            `json:"type,omitempty"`
	Credential *models.RequestExchangeCredential `json:"credential"`
	Timestamp  time.Time                         `json:"timestamp,omitempty"`
}

type CredentialKafkaRepository struct {
	client KafkaClient
}

func NewCredentialKafkaRepository(client KafkaClient) *CredentialKafkaRepository {
	return &CredentialKafkaRepository{
		client: client,
	}
}

func (c *CredentialKafkaRepository) UpdateCredential(payload *models.RequestExchangeCredential, token *oauth2.Token) (sended bool) {
	// not necessary returned value btw for more redeability
	payload = c.UpdateCredentialFromGoogle(payload, token)
	command := CredentialCommand{
		Credential: payload,
	}
	sended = c.PublishCommand(command, payload.ID)
	return sended
}

func (c *CredentialKafkaRepository) PublishCommand(credentialCommand CredentialCommand, key string) bool {
	command, err := json.Marshal(credentialCommand)
	if err != nil {
		log.Printf("ERROR | Cannot transform to JSON %v", err)
		return false
	}

	for i := 1; i < models.MaxAttempts; i++ {
		err = c.client.Produce("credentials.command", []byte(key), command)
		if err == nil {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to Broker, attempt %d: %v. Retrying in %v", i, err, waitTime)
		time.Sleep(waitTime)
	}

	return false
}

// TODO: it can be moved to another position
// not necessary returned value btw for more redeability
func (c *CredentialKafkaRepository) UpdateCredentialFromGoogle(exchangeCredential *models.RequestExchangeCredential, token *oauth2.Token) *models.RequestExchangeCredential {
	exchangeCredential.Data.Token = token.AccessToken
	exchangeCredential.Data.TokenRefresh = token.RefreshToken
	exchangeCredential.UpdatedAt.Time = time.Now().UTC()

	// token.expiry already set to 0
	exchangeCredential.ExpiresAt.Time = token.Expiry.UTC().Add(-models.TimeDriftForExpire * time.Second)
	return exchangeCredential
}
