package di

import (
	"actions_google/pkg/config"
	"actions_google/pkg/dimodel"
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/services"
	"actions_google/pkg/infra/brokerclient"
	"actions_google/pkg/infra/httpclient"
	"actions_google/pkg/infra/redisclient"
	"actions_google/pkg/infra/transform"
	"actions_google/pkg/interfaces/controllers"
)

// InitDependencies initializes and returns a pointer to a Dependencies struct
// containing all the necessary dependencies for the application.
//
// It sets up configurations for Zitadel, Kafka, and Clickhouse, and initializes
// authentication context, services, and controllers. It also creates clients
// for HTTP, Redis, and Kafka, and repositories for credentials and actions.
//
// The returned Dependencies struct includes:
// - AuthService: a pointer to the authentication service
// - AuthController: the authentication controller
// - ActionsController: the actions controller
//
// Returns:
// - *dimodel.Dependencies: a pointer to the initialized Dependencies struct
func InitDependencies() *dimodel.Dependencies {
	configZitadel := config.NewZitaldelEnvConfig()
	kafkaConfig := config.NewKafkaEnvConfig()
	clickhouseConfig := config.NewClickhouseEnvConfig()

	// init autentication
	authContext := controllers.NewAuthContext(configZitadel)
	authService := authContext.GetAuthService()
	authController := authContext.GetAuthController()

	credentialBrokerClient := brokerclient.NewBrokerClient(kafkaConfig)
	repoCredentialBroker := brokerclient.NewCredentialKafkaRepository(credentialBrokerClient)

	actionsHTTPClient := httpclient.NewClientImpl(models.TimeoutRequest)
	credentialHTTPClient := httpclient.NewClientImpl(models.TimeoutRequest)

	repoCredentialHTTP := httpclient.NewCredentialRepository(credentialHTTPClient, clickhouseConfig)
	actionsRedisClient := redisclient.NewRedisClient()
	actionsBrokerClient := brokerclient.NewBrokerClient(kafkaConfig)
	notionRepo := transform.NewActionsClient()
	repoActionsRedis := redisclient.NewActionsRepository(actionsRedisClient)
	repoActionsBroker := brokerclient.NewActionsKafkaRepository(actionsBrokerClient)
	actionsRepo := httpclient.NewActionsClientHTTP(actionsHTTPClient, clickhouseConfig)

	tokenAuth := services.NewTokenAuthImpl()
	sheetUtils := services.NewSheetUtilsImpl()

	actionsService := services.NewActionsService(
		repoActionsRedis,
		repoActionsBroker,
		actionsRepo,
		repoCredentialHTTP,
		repoCredentialBroker,
		notionRepo,
		tokenAuth,
		sheetUtils,
	)
	actionsController := controllers.NewActionsController(actionsService)

	return &dimodel.Dependencies{
		AuthService:       &authService,
		AuthController:    authController,
		ActionsController: actionsController,
	}
}
