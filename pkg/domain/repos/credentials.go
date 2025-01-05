package repos

import "actions_google/pkg/domain/models"

type CredentialHTTPRepository interface {
	GetCredentialByID(userID *string, credentialID *string, limitCount uint64) (*models.RequestExchangeCredential, error)
	GetAllCredentials(userID *string, limitCount uint64) (*[]models.RequestExchangeCredential, error)
}
