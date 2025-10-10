package jobs

import (
	"testing"

	slack "github.com/slack-go/slack"
)

func TestSendSlackMessage_DevMode(t *testing.T) {
	slackClient := &Slack{}
	slackClient.client = slack.New("slack_token")
	args := SlackArgs{
		Channel: "general",
		Message: "Hello from dev mode!",
		DevMode: true,
	}
	if err := slackClient.SendSlackMessage(args); err != nil {
		t.Fatalf("expected no error in dev mode, got: %v", err)
	}
}
