package jobs

import (
	"context"
	"testing"
)

func TestSendSlackMessage_DevMode(t *testing.T) {
	ctx := context.Background()
	args := SlackArgs{
		Channel: "general",
		Message: "Hello from dev mode!",
		DevMode: true,
	}
	if err := SendSlackMessage(ctx, args); err != nil {
		t.Fatalf("expected no error in dev mode, got: %v", err)
	}
}
