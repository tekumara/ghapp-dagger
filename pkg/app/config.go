package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/palantir/go-githubapp/githubapp"
)

func GhaConfig() (*githubapp.Config, error) {
	baseURL, exists := os.LookupEnv("GITHUB_BASE_URL")
	if !exists {
		return nil, fmt.Errorf("Missing env var GITHUB_BASE_URL.")
	}

	privateKey, err := ioutil.ReadFile(os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"))
	if err != nil {
		return nil, fmt.Errorf("Unable to load GitHub App private key from file: %s", os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"))
	}

	id, err := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Missing or non-integer env var GITHUB_APP_ID.")
	}

	secret, exists := os.LookupEnv("GITHUB_APP_WEBHOOK_SECRET")
	if !exists {
		return nil, fmt.Errorf("Missing env var GITHUB_APP_WEBHOOK_SECRET.")
	}

	app := &githubapp.Config{
		V3APIURL: baseURL,
		App: struct {
			IntegrationID int64  `yaml:"integration_id" json:"integrationId"`
			WebhookSecret string `yaml:"webhook_secret" json:"webhookSecret"`
			PrivateKey    string `yaml:"private_key" json:"privateKey"`
		}{
			IntegrationID: id,
			WebhookSecret: secret,
			PrivateKey:    string(privateKey),
		},
	}

	return app, nil
}
