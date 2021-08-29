package config

import "os"

const (
	apiGithubAccessToken = "GITHUB_ACCESS_TOKEN"
	LogLevel             = "info"
)

var (
	githubAccessToken = os.Getenv(apiGithubAccessToken)
)

func GetGithubAccessToken() string {
	return githubAccessToken
}

func IsProduction() bool {
	return os.Getenv("GO_ENV") == "production"
}
