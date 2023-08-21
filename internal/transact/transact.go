package transact

import (
	"service-auth-cff/internal/dbx"
	"service-auth-cff/pkg/transactions/loggedusers"

	"service-auth-cff/internal/logger"
)

func RegisterConfig(action string, description string, user string) {

}

func RegisterLogUsr(Event string, HostName string, IpRequest string, IpRemote string, UserId string) {
	conn := dbx.GetConnection()

	repoLoggedUsers := loggedusers.FactoryStorage(conn, nil, "")
	srvLoggedUsers := loggedusers.NewTxLoggedUserService(repoLoggedUsers, nil, "")

	_, _, err := srvLoggedUsers.CreateTxLoggedUser(Event, HostName, IpRequest, IpRemote, UserId)
	if err != nil {
		logger.Error.Println("", " - couldn't create loggedUsers :", err)
	}
}

func RegisterTrace(typeMessage string, fileName string, codeLine int, message string, transaction string) {

}
