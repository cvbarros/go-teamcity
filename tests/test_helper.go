package tests

import teamcity "github.com/cvbarros/go-teamcity-sdk"

var (
	client   *teamcity.Client
	username string
	password string
)

func initTest() *teamcity.Client {
	username = "admin"
	password = "admin"

	return teamcity.NewClient(username, password)
}

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	b := false
	return &b
}
