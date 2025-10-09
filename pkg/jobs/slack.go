// Package jobs provides river jobs for various integrations.
package jobs

import (
	"context"
	"errors"
	"fmt"

	slack "github.com/slack-go/slack"
)
// slackConversationsLimit is the maximum number of conversations to fetch from Slack.
const slackConversationsLimit = 1000

// errChannelNotFound is the error returned when a channel is not found.
var errChannelNotFound = errors.New("channel not found")

// errCouldNotFindChannel is the error returned when a channel could not be found.
var errCouldNotFindChannel = errors.New("could not find channel")

// errSlackTokenMissing is the error returned when the Slack token is missing.
var errSlackTokenMissing = errors.New("slack token is missing")

// SendSlackMessage sends a message to a Slack channel using a Slack App.
func SendSlackMessage(ctx context.Context, args SlackArgs) error {
	if args.DevMode {
		fmt.Printf("[DEV MODE] Would send to channel '%s': %s\n", args.Channel, args.Message)
		return nil
	}

	token, ok := ctx.Value(slackTokenKey).(string)
	if !ok || token == "" {
		return errSlackTokenMissing
	}

	client := slack.New(token)
	channelID := args.Channel

	// If channel is not an ID, try to resolve name
	if !isChannelID(args.Channel) {
		ch, err := findChannelByName(client, args.Channel)
		if err != nil {
			return errCouldNotFindChannel
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
    params := &slack.GetConversationsParameters{
        Types:           []string{"public_channel", "private_channel"},
        ExcludeArchived: true,
        Limit:           slackConversationsLimit,
    }
    channels, _, err := client.GetConversations(params)
    if err != nil {
        return nil, err
    }
    for _, ch := range channels {
        if ch.Name == name {
            return &ch, nil
        }
    }
    return nil, fmt.Errorf("channel '%s' not found: %w", name, errChannelNotFound)
}
