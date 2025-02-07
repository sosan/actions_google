package controllers

import (
	"actions_google/pkg/auth"
	"actions_google/pkg/config"
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/repos"
	"actions_google/pkg/domain/services"
	"actions_google/pkg/infra/httpclient"
	"actions_google/pkg/infra/redisclient"
	"actions_google/pkg/infra/tokenrepo"
	"actions_google/pkg/interfaces/middlewares"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService repos.AuthService
	once        sync.Once
	Config      config.ZitadelConfig
}

func NewAuthContext(cfg config.ZitadelConfig) *AuthController {
	return &AuthController{
		Config: cfg,
	}
}

func (ac *AuthController) GetAuthController() *AuthController {
	ac.once.Do(func() {
		zitadelClient := httpclient.NewZitadelClient(
			ac.Config.GetZitadelURI(),
			ac.Config.GetZitadelServiceUserID(),
			ac.Config.GetZitadelServiceUserKeyPrivate(),
			ac.Config.GetZitadelServiceUserKeyID(),
			ac.Config.GetZitadelProjectID(),
			ac.Config.GetZitadelKeyClientID(),
		)
		jwtGenerator := auth.NewJWTGenerator(auth.JWTGeneratorConfig{
			ServiceUser: auth.ServiceUserConfig{
				UserID:     ac.Config.GetZitadelServiceUserID(),
				PrivateKey: []byte(ac.Config.GetZitadelServiceUserKeyPrivate()),
				KeyID:      ac.Config.GetZitadelServiceUserKeyID(),
				ClientID:   ac.Config.GetZitadelServiceUserClientID(),
			},
			BackendApp: auth.BackendAppConfig{
				AppID:      ac.Config.GetZitadelBackendID(),
				PrivateKey: []byte(ac.Config.GetZitadelBackendKeyPrivate()),
				KeyID:      ac.Config.GetZitadelBackendKeyID(),
				ClientID:   ac.Config.GetZitadelBackendClientID(),
			},
			APIURL:    ac.Config.GetZitadelURI(),
			ProjectID: ac.Config.GetZitadelProjectID(),
			ClientID:  ac.Config.GetZitadelKeyClientID(),
		})
		redisClient := redisclient.NewRedisClient()
		tokenRepo := tokenrepo.NewTokenRepository(redisClient)
		authService := services.NewAuthService(jwtGenerator, zitadelClient, tokenRepo)

		// cache token logic
		cachedToken := authService.GetCachedActionUserAccessToken()
		// in dev state, not rotating service user acces token in servesless functions
		if ac.Config.GetEnv("ROTATE_SERVICE_USER_TOKEN", "n") == "y" {
			if cachedToken == nil {
				// Rotate token if it's expired or not found
				_, err := authService.GenerateAccessToken()
				if err != nil {
					_ = authService.GetCachedActionUserAccessToken()
				}
			}
		}
		ac.AuthService = authService
	})
	return ac
}

func (ac *AuthController) GetAuthService() repos.AuthService {
	return ac.GetAuthController().AuthService
}

func (ac *AuthController) VerifyUserToken(ctx *gin.Context) {
	userToken := ctx.Param("usertoken")
	isValid, isExpired := ac.AuthService.VerifyUserToken(userToken)
	if !isValid {
		ctx.JSON(http.StatusUnauthorized, middlewares.NewUnauthorizedError(models.AuthInvalid))
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"valid":   isValid,
		"expired": isExpired,
		"error":   "",
	})
}

func (ac *AuthController) VerifyUserTokenForMiddleware(ctx *gin.Context) {
	userToken := ctx.Param("usertoken")
	isValid, isExpired := ac.AuthService.VerifyUserToken(userToken)
	if !isValid {
		ctx.JSON(http.StatusUnauthorized, middlewares.NewUnauthorizedError(models.AuthInvalid))
		ctx.Abort()
		return
	}
	if isExpired {
		ctx.JSON(http.StatusUnauthorized, middlewares.NewUnauthorizedError(models.UserTokenExpired))
		ctx.Abort()
		return
	}
	ctx.Next()
}
