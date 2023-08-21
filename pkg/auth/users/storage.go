package users

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

type ServicesUserRepository interface {
	Create(m *User) error
	Update(m *User) error
	Delete(id string) error
	GetByID(id string) (*User, error)
	GetAll() ([]*User, error)
	GetByUsername(username string) (*User, error)
	GetUsersByIDs(ids []string) ([]*User, error)
	BlockUser(id string) error
	UnblockUser(id string) error
	LogoutUser(id string) error
	ChangePassword(id string, password string) error
	UpdatePasswordByUser(id string, password string) error
	GetByUsernameAndIdentificationNumber(username string, identificationNumber string) (*User, error)
	DeleteUserPasswordHistory(id string) error
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUserRepository {
	var s ServicesUserRepository
	engine := db.DriverName()
	switch engine {
	case SqlServer:
		return NewUserSqlServerRepository(db, user, txID)
	case Postgresql:
		return NewUserPsqlRepository(db, user, txID)
	case Oracle:
		return NewUserOrclRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
