package brokerclient

import (
	"encoding/json"
	// "fmt"
	"actions_google/pkg/common"
	"actions_google/pkg/domain/models"
	"log"
	"time"
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

func (c *CredentialKafkaRepository) UpdateCredential(payload *models.RequestExchangeCredential) (sended bool) {
	// payload := c.credentialToPayload(stateInfo, token, refresh, expire)
	// if payload == nil {
	// 	return false
	// }
	command := CredentialCommand{
		Credential: payload,
	}
	// key := fmt.Sprintf("credential_%s_%s_%s_%s", stateInfo.Sub, stateInfo.WorkflowID, stateInfo.NodeID, stateInfo.Type)
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
