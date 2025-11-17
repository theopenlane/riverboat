// Package jobs provides river jobs for various integrations.
package jobs

import (
	"context"
	"errors"
	"fmt"

	slack "github.com/slack-go/slack"
)

// Slack holds the address of client for sending a Slack message.
type Slack struct {
	client *slack.Client
}

// slackConversationsLimit is the maximum number of conversations to fetch from Slack.
const slackConversationsLimit = 1000

// errChannelNotFound is the error returned when a channel is not found.
var errChannelNotFound = errors.New("channel not found")

// errCouldNotFindChannel is the error returned when a channel could not be found.
var errCouldNotFindChannel = errors.New("could not find channel")

// SendSlackMessage sends a message to a Slack channel using a Slack App.
func (s *Slack) SendSlackMessage(ctx context.Context, args SlackArgs) error {
	if args.DevMode {
		fmt.Printf("[DEV MODE] Would send to channel '%s': %s\n", args.Channel, args.Message)
		return nil
	}

	channelID := args.Channel

	// If channel is not an ID, try to resolve name
	if !isChannelID(args.Channel) {
		ch, err := s.findChannelByName(args.Channel)
		if err != nil {
			return errCouldNotFindChannel
		}

		channelID = ch.ID
	}

	_, _, err := s.client.PostMessageContext(ctx, channelID, slack.MsgOptionText(args.Message, false))
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
func (s *Slack) findChannelByName(name string) (*slack.Channel, error) {
	params := &slack.GetConversationsParameters{
		Types:           []string{"public_channel", "private_channel"},
		ExcludeArchived: true,
		Limit:           slackConversationsLimit,
	}

	channels, _, err := s.client.GetConversations(params)
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
