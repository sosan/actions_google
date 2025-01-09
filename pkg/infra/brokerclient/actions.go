package brokerclient

import (
	"actions_google/pkg/common"
	"actions_google/pkg/domain/models"
	"encoding/json"
	"log"
	"time"
)

type ActionsKafkaRepository struct {
	client KafkaClient
}

const (
	CommandTypeCreate = "create"
	CommandTypeUpdate = "update"
	CommandTypeDelete = "delete"
	TopicName         = "actions.command"
)

type ActionsCommand struct {
	Actions   *models.RequestGoogleAction `json:"actions"`
	Type      string                      `json:"type,omitempty"`
	Timestamp time.Time                   `json:"timestamp,omitempty"`
}

func NewActionsKafkaRepository(client KafkaClient) *ActionsKafkaRepository {
	return &ActionsKafkaRepository{
		client: client,
	}
}

func (a *ActionsKafkaRepository) SendAction(newAction *models.RequestGoogleAction) (sended bool) {
	command := ActionsCommand{
		Actions:   newAction,
		Type:      CommandTypeUpdate,
		Timestamp: time.Now().UTC(),
	}
	sended = a.PublishCommand(command, newAction.ActionID)
	return sended
}

func (a *ActionsKafkaRepository) PublishCommand(payload ActionsCommand, key string) bool {
	command, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERROR | Cannot transform to JSON %v", err)
		return false
	}

	for i := 1; i < models.MaxAttempts; i++ {
		err = a.client.Produce(TopicName, []byte(key), command)
		if err == nil {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to Broker, attempt %d: %v. Retrying in %v", i, err, waitTime)
		time.After(waitTime)
	}

	return false
}
