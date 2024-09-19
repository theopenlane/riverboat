package jobs

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/providers/resend"
)

// EmailArgs for the email worker to process the job
type EmailArgs struct {
	// Message is the email message to send
	Message newman.EmailMessage `json:"message"`
}

// Kind satisfies the river.Job interface
func (EmailArgs) Kind() string { return "email" }

// EmailWorker is a worker to send emails using the resend email provider
// the config defaults to dev mode, which will write the email to a file using the mock provider
// a token is required to send emails using the actual resend provider
type EmailWorker struct {
	river.WorkerDefaults[EmailArgs]

	EmailConfig
}

// EmailConfig contains the configuration for the email worker
type EmailConfig struct {
	// DevMode is a flag to enable dev mode
	DevMode bool `koanf:"devMode" json:"devMode" jsonschema:"description=enable dev mode" default:"true"`
	// TestDir is the directory to use for dev mode
	TestDir string `koanf:"testDir" json:"testDir" jsonschema:"description=the directory to use for dev mode" default:"fixtures/email"`
	// Token is the token to use for the email provider
	Token string `koanf:"token" json:"token" jsonschema:"description=the token to use for the email provider"`
	// FromEmail is the email address to use as the sender
	FromEmail string `koanf:"fromEmail" json:"fromEmail" jsonschema:"required description=the email address to use as the sender" default:"no-reply@example.com"`
}

// validateEmailConfig validates the email configuration settings
func (w *EmailWorker) validateEmailConfig() error {
	if w.DevMode && w.TestDir == "" {
		return newMissingRequiredArg("test directory", EmailArgs{}.Kind())
	}

	if !w.DevMode && w.Token == "" {
		return newMissingRequiredArg("token", EmailArgs{}.Kind())
	}

	return nil
}

// Work satisfies the river.Worker interface for the email worker
// it sends an email using the resend email provider with the provided email message
func (w *EmailWorker) Work(ctx context.Context, job *river.Job[EmailArgs]) error {
	// validate the email configuration
	if err := w.validateEmailConfig(); err != nil {
		return err
	}

	log.Info().Strs("to", job.Args.Message.To).
		Str("subject", job.Args.Message.Subject).
		Msg("sending email")

	// set the options for the resend client
	opts := []resend.Option{}

	if w.DevMode {
		log.Debug().Str("directory", w.TestDir).Msg("running in dev mode")

		opts = append(opts, resend.WithDevMode(w.TestDir))
	}

	// if the from email is not set on the message, use the default from the worker config
	if job.Args.Message.From == "" {
		job.Args.Message.From = w.FromEmail
	}

	// create the resend client
	client, err := resend.New(w.Token, opts...)
	if err != nil {
		return err
	}

	return client.SendEmailWithContext(ctx, &job.Args.Message)
}
