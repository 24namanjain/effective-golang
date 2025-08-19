package slack

import (
	"context"
	"fmt"
	"time"

	githubslack "github.com/slack-go/slack"
	"slack-notifier/internal/events"
)

// Client wraps the slack-go client with a minimal API we need.
type Client struct {
	api     *githubslack.Client
	channel string
}

// NewClient constructs a new Client using the bot token and default channel.
func NewClient(botToken, channel string) *Client {
	return &Client{
		api:     githubslack.New(botToken),
		channel: channel,
	}
}

// SendMessage posts a plain text message to the configured channel.
func (c *Client) SendMessage(ctx context.Context, message string) error {
	_, _, err := c.api.PostMessageContext(ctx, c.channel, githubslack.MsgOptionText(message, false))
	return err
}

// SendEvent posts a rich-formatted message for an event. If event.Channel is empty, default channel is used.
func (c *Client) SendEvent(ctx context.Context, event *events.Event) error {
	channel := event.Channel
	if channel == "" {
		channel = c.channel
	}

	blocks := c.buildBlocks(event)
	_, _, err := c.api.PostMessageContext(ctx, channel, githubslack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("post slack message: %w", err)
	}
	return nil
}

// TestConnection validates the token by calling auth.test.
func (c *Client) TestConnection(ctx context.Context) error {
	_, err := c.api.AuthTestContext(ctx)
	return err
}

func (c *Client) buildBlocks(event *events.Event) []githubslack.Block {
	var blocks []githubslack.Block

	header := fmt.Sprintf("*%s* %s", event.Title, emojiFor(event.Severity))
	blocks = append(blocks, githubslack.NewSectionBlock(
		githubslack.NewTextBlockObject("mrkdwn", header, false, false), nil, nil,
	))

	if event.Message != "" {
		blocks = append(blocks, githubslack.NewSectionBlock(
			githubslack.NewTextBlockObject("mrkdwn", event.Message, false, false), nil, nil,
		))
	}

	ctxElems := []githubslack.MixedElement{
		githubslack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Type:* %s", event.Type), false, false),
		githubslack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Time:* %s", event.Timestamp.Format(time.RFC3339)), false, false),
		githubslack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*ID:* `%s`", event.ID), false, false),
	}
	if event.UserID != "" {
		ctxElems = append(ctxElems, githubslack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*User:* %s", event.UserID), false, false))
	}
	blocks = append(blocks, githubslack.NewContextBlock("", ctxElems...))

	if len(event.Metadata) > 0 {
		md := "*Details:*\n"
		for k, v := range event.Metadata {
			md += fmt.Sprintf("• %s: %v\n", k, v)
		}
		blocks = append(blocks, githubslack.NewSectionBlock(
			githubslack.NewTextBlockObject("mrkdwn", md, false, false), nil, nil,
		))
	}

	blocks = append(blocks, githubslack.NewDividerBlock())
	return blocks
}

func emojiFor(s events.Severity) string {
	switch s {
	case events.SeverityInfo:
		return "ℹ️"
	case events.SeverityWarning:
		return "⚠️"
	case events.SeverityError:
		return "❌"
	case events.SeverityCritical:
		return "🚨"
	default:
		return ""
	}
}
