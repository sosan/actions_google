package services

import (
	"fmt"
	"log"
	"actions_google/pkg/common"
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/repos"
	"time"

	"github.com/google/uuid"
)

type ActionsServiceImpl struct {
	redisRepo  repos.ActionsRedisRepoInterface
	brokerRepo repos.ActionsBrokerRepository
	httpRepo   repos.ActionsHTTPRepository
}

func NewActionsService(repoRedis repos.ActionsRedisRepoInterface, repoBroker repos.ActionsBrokerRepository, repoHTTP repos.ActionsHTTPRepository) repos.ActionsService {
	return &ActionsServiceImpl{
		redisRepo:  repoRedis,
		brokerRepo: repoBroker,
		httpRepo:   repoHTTP,
	}
}

func (a *ActionsServiceImpl) GetGoogleSheetByID(newAction models.RequestGoogleAction) (created bool, exist bool, action *models.ActionData) {
	now := time.Now().UTC().Format(models.LayoutTimestamp)
	exist, err := a.SetActionID(&newAction, &now)
	if err != nil || exist { //  in case that 10 loops cannot get new UUID just return because cannot get new uuid
		return false, exist, nil
	}

	for i := 1; i < models.MaxAttempts; i++ {
		created, exist = a.retriesCreateAction(&newAction, now)
		if !exist && created {
			return created, exist, nil
		}
		if exist {
			return false, true, nil
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("WARNING | Failed to create workflow, attempt %d:. Retrying in %v", i, waitTime)
		time.Sleep(waitTime)
	}
	log.Print("ERROR | Needs to add to Dead Letter. Cannot create workflow")
	// TODO: dead letter
	return false, false, nil
}

func (a *ActionsServiceImpl) ValidateActionGlobalUUID(field *string) (bool, error) {
	return a.redisRepo.ValidateActionGlobalUUID(field)
}

func (a *ActionsServiceImpl) retriesCreateAction(newAction *models.RequestGoogleAction, now string) (created bool, exist bool) {
	newAction.CreatedAt = now
	// created, exist, err := a.redisRepo.Create(newAction)
	// if err != nil {
	// 	log.Printf("ERROR | acquiring lock: %v", err)
	// 	return false, false
	// }
	// if exist || !created {
	// 	return created, exist
	// }
	// // if !createdRedis {
	// // 	return false, false
	// // }

	sended := a.brokerRepo.Create(newAction)
	if !sended {
		log.Printf("ERROR | Failed to publish action event %v", newAction)
		// a.redisRepo.Remove(newAction)
		return false, false
	}
	return created, exist
}

// TODO: maybe can make general function to create requestID and IDs
func (a *ActionsServiceImpl) generateActionID(now *string) (string, string) {
	actionID := uuid.New().String()
	if actionID == "" { // in case fail
		actionID = uuid.New().String()
	}
	requestID := fmt.Sprintf("%s_%s", actionID, *now)
	return actionID, requestID
}

func (a *ActionsServiceImpl) SetActionID(newAction *models.RequestGoogleAction, now *string) (exist bool, err error) {
	var actionID, requestID string

	for i := 0; i < models.MaxAttempts; i++ {
		actionID, requestID = a.generateActionID(now)
		exist, err = a.ValidateActionGlobalUUID(&actionID)
		if err != nil {
			exist = true
			return exist, err
		}
		if !exist {
			break
		}
		// not used time.sleep
	}

	// All attempts failed to find a unique UUID
	if exist {
		return exist, fmt.Errorf("all attempts to generate a unique UUID failed")
	}

	newAction.ActionID = actionID
	newAction.RequestID = requestID
	return exist, nil
}
