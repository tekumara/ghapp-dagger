package app

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/v45/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type CheckSuiteEventHandler struct {
	githubapp.ClientCreator
}

func (h *CheckSuiteEventHandler) Handles() []string {
	return []string{"check_suite"}
}

func (h *CheckSuiteEventHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.CheckSuiteEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse check suite event payload")
	}

	// TODO: structured logging
	zerolog.Ctx(ctx).Info().Msgf("check_suite: action %s owner %s repo %s headSHA %s\n",
		*event.Action, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckSuite.HeadSHA,
	)

	if *event.Action == "requested" || *event.Action == "rerequested" {
		installationID := githubapp.GetInstallationIDFromEvent(&event)
		client, err := h.NewInstallationClient(installationID)
		if err != nil {
			return err
		}
		return CreateCheckRun(ctx, client, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckSuite.HeadSHA)
	}

	return nil
}

type CheckRunEventHandler struct {
	githubapp.ClientCreator
	AppID int64
}

func (h *CheckRunEventHandler) Handles() []string {
	return []string{"check_run"}
}

func (h *CheckRunEventHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.CheckRunEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse check run event payload")
	}

	zerolog.Ctx(ctx).Info().Msgf("check_run: app %d action %s id %d owner %s repo %s headSHA %s\n",
		*event.CheckRun.App.ID, *event.Action, *event.CheckRun.ID,
		*event.Repo.Owner.Login, *event.Repo.Name, *event.CheckRun.HeadSHA,
	)

	// we receive checks runs created by other apps installed in the repo, so only process our ones
	if *event.CheckRun.App.ID == h.AppID {

		installationID := githubapp.GetInstallationIDFromEvent(&event)
		client, err := h.NewInstallationClient(installationID)
		if err != nil {
			return err
		}

		// the user has pressed "Re-run"
		if *event.Action == "rerequested" {
			// TODO: handle errors
			if err := CreateCheckRun(ctx, client, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckRun.HeadSHA); err != nil {
				//TODO: structured logging with context of request
				zerolog.Ctx(ctx).Error().Msgf("%v", err)
			}
		}

		// run the check
		if *event.Action == "created" {
			//TODO: structured logging with context of request
			zerolog.Ctx(ctx).Info().Msg("check_run: execute check!")
			if err := ExecuteCheck(ctx, client, &event); err != nil {
				//TODO: structured logging with context of request
				zerolog.Ctx(ctx).Error().Msgf("%v", err)
			}
		}

	}

	return nil
}
