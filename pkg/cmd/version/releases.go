package version

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/google/go-github/v41/github"
	"golang.org/x/net/context"
)

func LatestRelease() (string, error) {
	httpClient := http.Client{}
	defer httpClient.CloseIdleConnections()

	client := github.NewClient(&httpClient)
	latestRelease, _, err := client.Repositories.
		GetLatestRelease(context.Background(), "bartoszmajsak", "template-golang")
	if err != nil {
		return "", errors.Wrap(err, "unable to determine latest released version")
	}

	return *latestRelease.Name, nil
}

func IsLatestRelease(releaseVersion string) bool {
	latestRelease, _ := LatestRelease()

	return releaseVersion == latestRelease
}

