package main

import (
	"service-auth-cff/api"

	"service-auth-cff/internal/env"
)

func main() {
	c := env.NewConfiguration()
	api.Start(c.App.Port, c.App.ServiceName, c.App.LoggerHttp, c.App.AllowedDomains)
}
