package password

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

type ServicesPasswordRepository interface {
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesPasswordRepository {
	var s ServicesPasswordRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		fallthrough
	case Postgresql:
		fallthrough
	case Oracle:
		fallthrough
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
