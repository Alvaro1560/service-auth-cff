package auth

import (
	"github.com/jmoiron/sqlx"
	"service-auth-cff/internal/models"
	"service-auth-cff/pkg/auth/users"
	"service-auth-cff/pkg/auth/verification_email"
)

type ServerAuth struct {
	//SrvUsers users.PortsServerUser
	SrvVerificationEmail verification_email.PortsServerVerificationEmail
}

func NewServerAuth(db *sqlx.DB, usr *models.User, txID string) *ServerAuth {
	repoUsers := users.FactoryStorage(db, usr, txID)
	_ = users.NewUserService(repoUsers, usr, txID)

	repoVerificationEmail := verification_email.FactoryStorage(db, usr, txID)
	srvVerificationEmail := verification_email.NewVerificationEmailService(repoVerificationEmail, usr, txID)

	return &ServerAuth{
		//SrvUsers: srvUsers,
		SrvVerificationEmail: srvVerificationEmail,
	}

}
