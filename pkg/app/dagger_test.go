package app

import (
	"context"
	"os"
	"testing"
)

func TestExecDagger(t *testing.T) {
	ctx := context.Background()

	repoUrl := os.Getenv("GITHUB_REPO_URL")
	ref := os.Getenv("GITHUB_REF")
	token := os.Getenv("GITHUB_TOKEN")
	ExecDagger(ctx, repoUrl, ref, token, func(text string) error {
		// TODO: logging here
		return nil
	})
}
