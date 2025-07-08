package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/providers/resend"
)

const (
	maxEmailAttempts         = 2
	emailJobSnoozeDuration   = time.Second * 30
	emailRetryPolicyDuration = time.Minute
)

// EmailArgs for the email worker to process the job
type EmailArgs struct {
	// Message is the email message to send
	Message newman.EmailMessage `json:"message"`
}

// Kind satisfies the river.Job interface
func (EmailArgs) Kind() string { return "email" }

// InsertOpts provides the default configuration when processing this job.
// Here we want to retry sending an email a maxium of 3 times
// This can be overridden on inserting the job
func (EmailArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: maxEmailAttempts}
}

// EmailWorker is a worker to send emails using the resend email provider
// the config defaults to dev mode, which will write the email to a file using the mock provider
// a token is required to send emails using the actual resend provider
type EmailWorker struct {
	river.WorkerDefaults[EmailArgs]

	Config EmailConfig `koanf:"config" json:"config" jsonschema:"description=the email configuration"`

	client newman.EmailSender
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
	if w.Config.DevMode && w.Config.TestDir == "" {
		return newMissingRequiredArg("test directory", EmailArgs{}.Kind())
	}

	if !w.Config.DevMode && w.Config.Token == "" {
		return newMissingRequiredArg("token", EmailArgs{}.Kind())
	}

	// create the resend client only one
	if w.client == nil {
		// set the options for the resend client
		opts := []resend.Option{}

		if w.Config.DevMode {
			log.Debug().Str("directory", w.Config.TestDir).Msg("running in dev mode")

			opts = append(opts, resend.WithDevMode(w.Config.TestDir))
		}

		client, err := resend.New(w.Config.Token, opts...)
		if err != nil {
			return err
		}

		w.client = client
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

	// if the from email is not set on the message, use the default from the worker config
	if job.Args.Message.From == "" {
		job.Args.Message.From = w.Config.FromEmail
	}

	err := w.client.SendEmailWithContext(ctx, &job.Args.Message)
	if newman.IsRetryableError(err) {
		return river.JobSnooze(emailJobSnoozeDuration)
	}

	return err
}

// NextRetry always schedules the next retry for 30 seconds from now.
// In the case where we run into another error while processing the delivery of the email
// wait another 1 minute before retrying inside of almost immediately
//
// This might allow the email provider or others recover from the failure - if needed as against
// hammering the next request immediately
func (w *EmailWorker) NextRetry(_ *river.Job[EmailArgs]) time.Time {
	return time.Now().Add(emailRetryPolicyDuration)
}
