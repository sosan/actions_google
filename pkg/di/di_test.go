package di

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	// "time"

	"actions_google/mocks"
	config "actions_google/pkg/config"
	"actions_google/pkg/domain/services"
	tokenrepo "actions_google/pkg/infra/tokenrepo"
	controllers "actions_google/pkg/interfaces/controllers"
	tests "actions_google/tests"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	// "github.com/stretchr/testify/mock"
)

func TestVerifyUserToken_Success(t *testing.T) {
	// Arrange
	mockAuthService := new(mocks.AuthService)
	mockAuthService.On("VerifyUserToken", "valid-token").Return(true, false)

	controller := &controllers.AuthController{
		AuthService: mockAuthService,
	}

	router := gin.Default()
	router.GET("/verify/:usertoken", controller.VerifyUserToken)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/verify/valid-token", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	mockAuthService.AssertExpectations(t)
}

func TestInitDependencies(t *testing.T) {
	config.SetEnv("VAULT_URI", "redis://mock-redis-uri")
	config.SetEnv("TEST", "TEST")
	configZitadel := config.NewZitaldelEnvConfig()
	authController := controllers.NewAuthContext(configZitadel).GetAuthController()
	assert.NotNil(t, authController)
}

func TestInitDependencies_Failure(t *testing.T) {
	// Arrange
	mockAuthContext := new(mocks.TokenAuth)
	mockAuthContext.On("GetAuthService").Return(nil)

	// Act
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to missing AuthService")
		}
	}()
	InitDependencies()

	// Assert
	mockAuthContext.AssertExpectations(t)
}

func TestAuthenticate_Success(t *testing.T) {
	// Arrange
	mockTokenRepo := mocks.NewTokenRepository(t)
	mockTokenRepo.On("GetToken").Return(&tokenrepo.Token{AccessToken: tests.StringPtr("test-token")}, nil)

	mockJWTGenerator := mocks.NewJWTGenerator(t)
	// mockJWTGenerator.On("GenerateActionUserAssertionJWT", mock.Anything).Return("mock-jwt", nil)
	// mockJWTGenerator.On("GenerateAppInstrospectJWT", mock.Anything).Return("mock-app-jwt", nil)

	mockZitadelClient := mocks.NewZitadelClient(t)
	// mockZitadelClient.On("GenerateActionUserAccessToken", "mock-jwt").Return(tests.StringPtr("new-token"), 0, nil)
	// mockZitadelClient.On("ValidateUserToken", "valid-token", "mock-app-jwt").Return(true, 0, nil)
	// mockZitadelClient.On("ValidateActionUserAccessToken", tests.StringPtr("test-token"), tests.StringPtr("mock-app-jwt")).Return(true, nil)

	authService := services.NewAuthService(
		mockJWTGenerator,
		mockZitadelClient,
		mockTokenRepo,
	)

	result, err := authService.VerifyActionUserToken("test-token")

	assert.NoError(t, err)
	assert.True(t, result)
	mockTokenRepo.AssertExpectations(t)
	mockJWTGenerator.AssertExpectations(t)
	mockZitadelClient.AssertExpectations(t)
}

func TestVerifyUserToken_SuccessMocked(t *testing.T) {
	// Arrange
	mockJWTGenerator := new(mocks.JWTGenerator)
	mockZitadelClient := new(mocks.ZitadelClient)
	mockTokenRepo := new(mocks.TokenRepository)

	authService := services.NewAuthService(mockJWTGenerator, mockZitadelClient, mockTokenRepo)

	mockJWTGenerator.On("GenerateAppInstrospectJWT", mock.Anything).Return("mock-app-jwt", nil)
	mockZitadelClient.On("ValidateUserToken", "valid-user-token", "mock-app-jwt").Return(true, time.Now().UTC().Add(30000).Unix(), nil)

	// Act
	isValid, isExpired := authService.VerifyUserToken("valid-user-token")

	// Assert
	assert.True(t, isValid)
	assert.False(t, isExpired)

	mockJWTGenerator.AssertExpectations(t)
	mockZitadelClient.AssertExpectations(t)
}

func TestVerifyUserToken_FailedMocked(t *testing.T) {
	// Arrange
	mockJWTGenerator := new(mocks.JWTGenerator)
	mockZitadelClient := new(mocks.ZitadelClient)
	mockTokenRepo := new(mocks.TokenRepository)

	authService := services.NewAuthService(mockJWTGenerator, mockZitadelClient, mockTokenRepo)

	mockJWTGenerator.On("GenerateAppInstrospectJWT", mock.Anything).Return("mock-app-jwt", nil)
	mockZitadelClient.On("ValidateUserToken", "valid-user-token-invalid", "mock-app-jwt").Return(false, int64(3000), nil)

	// Act
	isValid, isExpired := authService.VerifyUserToken("valid-user-token-invalid")

	// Assert
	assert.False(t, isValid)
	assert.True(t, isExpired)

	mockJWTGenerator.AssertExpectations(t)
	mockZitadelClient.AssertExpectations(t)
}

func TestVerifyUserToken_ErrorFailedMocked(t *testing.T) {
	// Arrange
	mockJWTGenerator := new(mocks.JWTGenerator)
	mockZitadelClient := new(mocks.ZitadelClient)
	mockTokenRepo := new(mocks.TokenRepository)

	authService := services.NewAuthService(mockJWTGenerator, mockZitadelClient, mockTokenRepo)

	mockJWTGenerator.On("GenerateAppInstrospectJWT", mock.Anything).Return("mock-app-jwt", nil)
	mockZitadelClient.On("ValidateUserToken", "valid-user-token-invalid", "mock-app-jwt").
		Return(
			true,
			time.Now().UTC().Add(30000).Unix(),
			fmt.Errorf("error invalid"),
		)

	// Act
	isValid, isExpired := authService.VerifyUserToken("valid-user-token-invalid")

	// Assert
	assert.False(t, isValid)
	assert.True(t, isExpired)

	mockJWTGenerator.AssertExpectations(t)
	mockZitadelClient.AssertExpectations(t)
}
