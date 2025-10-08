// Package jobs provides river jobs for various integrations.
package jobs

import (
	"context"
	"fmt"
	"os"

	slack "github.com/slack-go/slack"
)


// SendSlackMessage sends a message to a Slack channel using a Slack App.
func SendSlackMessage(ctx context.Context, args SlackArgs) error {
	if args.DevMode {
		fmt.Printf("[DEV MODE] Would send to channel '%s': %s\n", args.Channel, args.Message)
		return nil
	}

	token, ok := ctx.Value("slack_token").(string)
	if !ok || token == "" {
		return fmt.Errorf("Slack token not provided in context")
	}

	client := slack.New(token)
	channelID := args.Channel

	// If channel is not an ID, try to resolve name
	if !isChannelID(args.Channel) {
		ch, err := findChannelByName(client, args.Channel)
		if err != nil {
			return fmt.Errorf("could not find channel '%s': %w", args.Channel, err)
		}
		channelID = ch.ID
	}

	_, _, err := client.PostMessage(channelID, slack.MsgOptionText(args.Message, false))
	if err != nil {
		return fmt.Errorf("failed to send Slack message: %w", err)
	}
	return nil
}

// isChannelID checks if the string looks like a Slack channel ID.
func isChannelID(s string) bool {
	return len(s) > 0 && (s[0] == 'C' || s[0] == 'G')
}

// findChannelByName looks up a channel by name.
func findChannelByName(client *slack.Client, name string) (*slack.Channel, error) {
	channels, err := client.GetChannels(false)
	if err != nil {
		return nil, err
	}
	for _, ch := range channels {
		if ch.Name == name {
			return &ch, nil
		}
	}
	return nil, fmt.Errorf("channel '%s' not found", name)
}
