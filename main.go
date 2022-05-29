package main

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	"github.com/swinton/go-probot/probot"
)

func main() {
	probot.HandleEvent("check_suite", func(ctx *probot.Context) error {
		event := ctx.Payload.(*github.CheckSuiteEvent)

		log.Printf("check_suite: action %s owner %s repo %s headSHA %s\n",
			*event.Action, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckSuite.HeadSHA)

		if *event.Action == "requested" || *event.Action == "rerequested" {
			createCheckRun(ctx.GitHub, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckSuite.HeadSHA)
		}

		return nil
	})

	probot.HandleEvent("check_run", func(ctx *probot.Context) error {
		event := ctx.Payload.(*github.CheckRunEvent)

		log.Printf("check_run: action %s owner %s repo %s headSHA %s\n",
			*event.Action, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckRun.HeadSHA)

		log.Printf("check_run: %s",jsonDump(event))

		// if *event.Action == "requested" || *event.Action == "rerequested" {
		// 	createCheckRun(ctx.GitHub, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckSuite.HeadSHA)
		// }

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
