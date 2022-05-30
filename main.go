package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/swinton/go-probot/probot"
)

func main() {
	// TODO: replace with https://github.com/palantir/go-githubapp because it has caches for client creation,
	// async webhook request handling, conditional requests on the client to avoid rate limiting, logging & metrics
	app := probot.NewApp()

	probot.HandleEvent("check_suite", func(ctx *probot.Context) error {
		event := ctx.Payload.(*github.CheckSuiteEvent)

		// TODO: structured logging
		log.Printf("check_suite: action %s owner %s repo %s headSHA %s\n",
			*event.Action, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckSuite.HeadSHA,
		)

		if *event.Action == "requested" || *event.Action == "rerequested" {
			createCheckRun(ctx.GitHub, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckSuite.HeadSHA)
		}

		return nil
	})

	probot.HandleEvent("check_run", func(ctx *probot.Context) error {
		event := ctx.Payload.(*github.CheckRunEvent)

		log.Printf("check_run: app %d action %s id %d owner %s repo %s headSHA %s\n",
			*event.CheckRun.App.ID, *event.Action, *event.CheckRun.ID,
			*event.Repo.Owner.Login, *event.Repo.Name, *event.CheckRun.HeadSHA,
		)

		// we receive checks runs created by other apps installed in the repo, so only process our ones
		if *event.CheckRun.App.ID == app.ID {

			// the user has pressed "Re-run"
			if *event.Action == "rerequested" {
				// TODO: handle errors
				createCheckRun(ctx.GitHub, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckRun.HeadSHA)
			}

			// run the check
			if *event.Action == "created" {
				log.Println("check_run: execute check!")
				if err := executeCheck(ctx.GitHub, event); err != nil {
					//TODO: structured logging with context of request
					log.Printf("ERROR: %v\n", err)
				}
			}

		}

		return nil
	})

	// start server (default port 8000)
	probot.Start()
}

func createCheckRun(ghClient *github.Client, owner, repo, headSHA string) {
	// TODO: check PRs from a fork are ignored
	// TODO: handle errors
	ghClient.Checks.CreateCheckRun(context.TODO(), owner, repo, github.CreateCheckRunOptions{
		Name:    "Demo Check",
		HeadSHA: headSHA,
	})
}

func executeCheck(ghClient *github.Client, event *github.CheckRunEvent) error {

	owner := *event.Repo.Owner.Login
	repo := *event.Repo.Name
	checkRunID := *event.CheckRun.ID
	ref := *event.CheckRun.HeadSHA
	repoUrl := *event.Repo.CloneURL
	appID := *event.CheckRun.App.ID
	installationID := *event.Installation.ID

	checkName := "Dagger"
	ctx := context.Background()

	updateCheckRunOutput := func(text string) error {
		_, _, err := ghClient.Checks.UpdateCheckRun(ctx, owner, repo, checkRunID, github.UpdateCheckRunOptions{
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
	ghAppClient, err := appAuthedClient(os.Getenv("GITHUB_BASE_URL"), os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"), appID)
	if err != nil {
		return err
	}

	installationToken, _, err := ghAppClient.Apps.CreateInstallationToken(ctx, installationID)
	if err != nil {
		// update run to complete with success or failure
		if _, _, err := ghClient.Checks.UpdateCheckRun(ctx, owner, repo, checkRunID, github.UpdateCheckRunOptions{
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
	output, execErr := execDagger(ctx, repoUrl, ref, token, updateCheckRunOutput)

	var conclusion string
	if execErr == nil {
		conclusion = "success"
	} else {
		conclusion = "failure"
		// TODO: return execErr?
		log.Printf("ERROR: execErr %+v", execErr)
	}

	// update run to complete with success or failure
	if _, _, err := ghClient.Checks.UpdateCheckRun(ctx, owner, repo, checkRunID, github.UpdateCheckRunOptions{
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
