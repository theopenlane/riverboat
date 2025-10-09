package jobs

import (
	"context"

	"github.com/riverqueue/river"
)

// contextKey is a custom type for context keys.
type contextKey string

const slackTokenKey contextKey = "slack_token"


// SlackArgs holds the arguments for a Slack message job.
type SlackArgs struct {
	// Channel can be a name or ID
	Channel string `json:"channel"`
	// Message is the text to send
	Message string `json:"message"`
	// DevMode mocks the request
	DevMode bool   `json:"dev_mode"`
}

// Kind returns the job kind for SlackArgs.
func (SlackArgs) Kind() string { return "slack" }

// SlackWorker sends messages to Slack.
type SlackWorker struct {
	river.WorkerDefaults[SlackArgs]

	Config SlackConfig `koanf:"config" json:"config" jsonschema:"description=the slack configuration"`
}
// Work satisfies the river.Worker interface for the Slack worker
// It sends a Slack message using the Slack App
func (w *SlackWorker) Work(ctx context.Context, job *river.Job[SlackArgs]) error {
	args := job.Args
	if !args.DevMode {
		args.DevMode = w.Config.DevMode
	}
	if w.Config.Token == "" && !args.DevMode {
		return newMissingRequiredArg("token", job.Args.Kind())
	}
	ctx = context.WithValue(ctx, slackTokenKey, w.Config.Token)
	return SendSlackMessage(ctx, args)
}

// SlackConfig configures the Slack worker.
type SlackConfig struct {
	// Enabled tells the server whether or not to register the worker, if disabled jobs of this type cannot be inserted
	Enabled bool   `koanf:"enabled" json:"enabled" jsonschema:"description=enable or disable the slack worker" default:"false"`
	// DevMode mocks the request, this can be overwritten on the individual job args
	DevMode bool   `koanf:"devMode" json:"devMode" jsonschema:"description=enable dev mode" default:"true"`
	// Token is the slack API token to connect to the slack instance
	Token   string `koanf:"token" json:"token" jsonschema:"description=the token to use for the slack app"`
}
