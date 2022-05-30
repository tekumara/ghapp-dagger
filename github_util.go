package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/github"
)

// auth as a github app see https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#authenticating-as-a-github-app
// needed to create a installation token
func appAuthedClient(baseURL string, privateKeyPath string, appID int64) (*github.Client, error) {
	privatePem, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pem: %v", err)
	}

	itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, appID, privatePem)
	if err != nil {
		return nil, fmt.Errorf("failed to create app transport: %v", err)
	}
	itr.BaseURL = baseURL

	//create git client with app transport
	client, err := github.NewEnterpriseClient(
		baseURL,
		baseURL,
		&http.Client{
			Transport: itr,
			Timeout:   time.Second * 30,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to create git client for app: %v", err)
	}

	return client, nil
}
