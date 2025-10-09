# Slack Job

This job sends a message to a Slack channel using a Slack App (not legacy webhook).

## Arguments
- `channel`: Channel name (e.g. "general") or channel ID (e.g. "C12345678")
- `message`: The message text
- `dev_mode`: If true, mocks the request (no token needed)

## Environment
- Requires `SLACK_BOT_TOKEN` to be set in the environment (unless in dev mode)

## Example Usage

```
river run slack --channel general --message "Hello, world!" --dev_mode true
```

## Implementation
- Uses [slack-go/slack](https://github.com/slack-go/slack)
- Resolves channel name to ID if needed
- See `slack.go` for details
