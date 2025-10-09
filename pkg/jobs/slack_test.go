package jobs

import (
	"context"
	"testing"

	slack "github.com/slack-go/slack"
)

func TestSendSlackMessage_DevMode(t *testing.T) {
	ctx := context.Background()
	slackClient := &Slack{}
	slackClient.client = slack.New("slack_token")
	args := SlackArgs{
		Channel: "general",
		Message: "Hello from dev mode!",
		DevMode: true,
	}
	if err := slackClient.SendSlackMessage(ctx, args); err != nil {
		t.Fatalf("expected no error in dev mode, got: %v", err)
	}
}
