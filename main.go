package main

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	"github.com/swinton/go-probot/probot"
)

func main() {
	// TODO: replace with https://github.com/palantir/go-githubapp because it has caches for client creation
	// & conditional requests to avoid rate limiting plus logging & metrics
	app := probot.NewApp()

	probot.HandleEvent("check_suite", func(ctx *probot.Context) error {
		event := ctx.Payload.(*github.CheckSuiteEvent)

		// TODO: structured logging
		log.Printf("check_suite: action %s owner %s repo %s headSHA %s\n",
			*event.Action, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckSuite.HeadSHA)

		if *event.Action == "requested" || *event.Action == "rerequested" {
			createCheckRun(ctx.GitHub, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckSuite.HeadSHA)
		}

		return nil
	})

	probot.HandleEvent("check_run", func(ctx *probot.Context) error {
		event := ctx.Payload.(*github.CheckRunEvent)

		log.Printf("check_run: app %d action %s owner %s repo %s headSHA %s\n",
			event.CheckRun.App.GetID(), *event.Action, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckRun.HeadSHA)

		// we receive checks runs created by other apps installed in the repo, so only process our ones
		if event.CheckRun.App.GetID() == app.ID {

			// the user has pressed "Re-run"
			if *event.Action == "rerequested" {
				createCheckRun(ctx.GitHub, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckRun.HeadSHA)
			}

			// run the check
			if *event.Action == "created" {
				log.Println("check_run: execute check!")
				executeCheck(ctx.GitHub, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckRun.ID)
			}

		}

		return nil
	})

	// start server (default port 8000)
	probot.Start()
}

func createCheckRun(ghClient *github.Client, owner, repo, headSHA string) {
	// TODO: check PRs from a fork are ignored
	ghClient.Checks.CreateCheckRun(context.TODO(), owner, repo, github.CreateCheckRunOptions{
		Name:    "Demo Check",
		HeadSHA: headSHA,
		// Conclusion: github.String("failure"),
		// Status:     github.String("completed"),
		// Output: &github.CheckRunOutput{
		// 	Title:   github.String("Config check failed"),
		// 	Summary: github.String(fmt.Sprintf(errString, event.ExecutionContext.ID, actionName)),
		// },
		// ExternalID: github.String(fmt.Sprintf("%s/%s", event.ExecutionContext.ID, actionName)),
	})
}

func executeCheck(ghClient *github.Client, owner, repo string, checkRunID int64) {
	ghClient.Checks.UpdateCheckRun(context.TODO(), owner, repo, checkRunID, github.UpdateCheckRunOptions{
		Name:   "Demo Check",
		Status: github.String("in_progress"),
	})

	ghClient.Checks.UpdateCheckRun(context.TODO(), owner, repo, checkRunID, github.UpdateCheckRunOptions{
		Name:       "Demo Check",
		Status:     github.String("completed"),
		Conclusion: github.String("success"),
	})
}
