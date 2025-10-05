package jobs

import (
	"context"
	"github.com/riverqueue/river"
)

// SlackArgs for the Slack worker to process the job
// Channel can be a name or ID
// Message is the text to send
// DevMode mocks the request
type SlackArgs struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
	DevMode bool   `json:"dev_mode"`
}

// river.Job interface implementation
func (SlackArgs) Kind() string { return "slack" }

// SlackWorker is a worker to send Slack messages
// Config contains the configuration for the Slack worker

type SlackWorker struct {
	river.WorkerDefaults[SlackArgs]

	Config SlackConfig `koanf:"config" json:"config" jsonschema:"description=the slack configuration"`
}
	// Work satisfies the river.Worker interface for the Slack worker
	// It sends a Slack message using the Slack App
	func (w *SlackWorker) Work(ctx context.Context, job *river.Job[SlackArgs]) error {
		if w.Config.DevMode || job.Args.DevMode {
		
			println("[DEV MODE] Would send to channel '", job.Args.Channel, "': ", job.Args.Message)
			return nil
		}

		token := w.Config.Token
		if token == "" {
			return newMissingRequiredArg("token", job.Args.Kind())
		}

		client := slack.New(token)
		channelID := job.Args.Channel

		// If channel is not an ID, try to resolve name
		if !isChannelID(job.Args.Channel) {
			ch, err := findChannelByName(client, job.Args.Channel)
			if err != nil {
				return err
			}
			channelID = ch.ID
		}

		_, _, err := client.PostMessage(channelID, slack.MsgOptionText(job.Args.Message, false))
		if err != nil {
			return err
		}
		return nil
	}

// SlackConfig contains the configuration for the Slack worker
// DevMode mocks the request
// Token is the Slack bot token
//
type SlackConfig struct {
	Enabled bool   `koanf:"enabled" json:"enabled" jsonschema:"description=enable or disable the slack worker" default:"false"`
	DevMode bool   `koanf:"devMode" json:"devMode" jsonschema:"description=enable dev mode" default:"true"`
	Token   string `koanf:"token" json:"token" jsonschema:"description=the token to use for the slack app"`
}
