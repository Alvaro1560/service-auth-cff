package verification_email

import (
	"github.com/jmoiron/sqlx"
	"service-auth-cff/internal/logger"
	"service-auth-cff/internal/models"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServicesVerificationEmailRepository interface {
	create(m *VerificationEmail) error
	update(m *VerificationEmail) error
	delete(id int64) error
	getByID(id int64) (*VerificationEmail, error)
	getAll() ([]*VerificationEmail, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesVerificationEmailRepository {
	var s ServicesVerificationEmailRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return newVerificationEmailSqlServerRepository(db, user, txID)
	case Postgresql:
		return newVerificationEmailPsqlRepository(db, user, txID)
	case Oracle:
		return newVerificationEmailOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
