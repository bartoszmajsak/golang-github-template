package version

import (
	"net/http"

	"github.com/bartoszmajsak/template-golang/version"

	"github.com/google/go-github/github"
	"golang.org/x/net/context"
)

func LatestRelease() (string, error) {
	httpClient := http.Client{}
	defer httpClient.CloseIdleConnections()

	client := github.NewClient(&httpClient)
	latestRelease, _, err := client.Repositories.
		GetLatestRelease(context.Background(), "bartoszmajsak", "template-golang")
	if err != nil {
		return "", err
	}
	return *latestRelease.Name, nil
}

func IsLatestRelease(latestRelease string) bool {
	return latestRelease == version.Version
}
