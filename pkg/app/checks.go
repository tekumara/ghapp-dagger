package app

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v45/github"
	"github.com/rs/zerolog"
)

func CreateCheckRun(ctx context.Context, client *github.Client, owner, repo, headSHA string) error {
	// TODO: check PRs from a fork are ignored
	_, _, err := client.Checks.CreateCheckRun(ctx, owner, repo, github.CreateCheckRunOptions{
		Name:    "Demo Check",
		HeadSHA: headSHA,
	})
	return err
}

func ExecuteCheck(ctx context.Context, client *github.Client, event *github.CheckRunEvent) error {

	owner := *event.Repo.Owner.Login
	repo := *event.Repo.Name
	checkRunID := *event.CheckRun.ID
	ref := *event.CheckRun.HeadSHA
	repoUrl := *event.Repo.CloneURL
	appID := *event.CheckRun.App.ID
	installationID := *event.Installation.ID

	checkName := "Dagger"

	// TODO: trim to 65535 chars
	updateCheckRunOutput := func(text string) error {
		_, _, err := client.Checks.UpdateCheckRun(ctx, owner, repo, checkRunID, github.UpdateCheckRunOptions{
			Name:   checkName,
			Status: github.String("in_progress"),
			Output: &github.CheckRunOutput{
				Title:   github.String(checkName),
				Summary: github.String("In progress ..."),
				Text:    github.String(text),
			},
		})
		return err
	}

	// update run to in progress (no output yet)
	if err := updateCheckRunOutput("....."); err != nil {
		return err
	}

	// TODO: move this into initial setup
	ghAppClient, err := AppAuthedClient(os.Getenv("GITHUB_BASE_URL"), os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"), appID)
	if err != nil {
		return err
	}

	installations, _, err := ghAppClient.Apps.ListInstallations(context.Background(), &github.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list installations: %v", err)
	}
	zerolog.Ctx(ctx).Info().Msgf("Installations: %+v", JsonDump(installations))

	// create installation token valid for 1 hour to use for cloning the repo
	installationToken, _, err := ghAppClient.Apps.CreateInstallationToken(ctx, installationID, &github.InstallationTokenOptions{})
	if err != nil {
		if _, _, err := client.Checks.UpdateCheckRun(ctx, owner, repo, checkRunID, github.UpdateCheckRunOptions{
			Name:       checkName,
			Status:     github.String("completed"),
			Conclusion: github.String("failure"),
			Output: &github.CheckRunOutput{
				Title:   github.String(checkName),
				Summary: github.String("Completed with failure"),
				Text:    github.String(err.Error()),
			},
		}); err != nil {
			return err
		}
		return fmt.Errorf("failed to create installation token: %v", err)
	}

	token := *installationToken.Token

	// execute dagger
	output, execErr := ExecDagger(ctx, repoUrl, ref, token, updateCheckRunOutput)

	var conclusion string
	if execErr == nil {
		conclusion = "success"
	} else {
		conclusion = "failure"
		// TODO: return execErr?
		zerolog.Ctx(ctx).Info().Msgf("ERROR: execErr %+v", execErr)
	}

	// update run to complete with success or failure
	// TODO: trim to 65535 chars
	if _, _, err := client.Checks.UpdateCheckRun(ctx, owner, repo, checkRunID, github.UpdateCheckRunOptions{
		Name:       checkName,
		Status:     github.String("completed"),
		Conclusion: github.String(conclusion),
		Output: &github.CheckRunOutput{
			Title:   github.String(checkName),
			Summary: github.String(fmt.Sprintf("Completed with %s", conclusion)),
			Text:    github.String(output),
		},
	}); err != nil {
		return err
	}

	return nil
}
