package testutil

import teamcity "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"

var (
	Client   *teamcity.Client
	username string
	password string
)

func InitTest() *teamcity.Client {
	username = "admin"
	password = "admin"
	Client = teamcity.NewClient(username, password)
	return Client
}
