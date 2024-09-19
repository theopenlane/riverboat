package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/theopenlane/newman"

	"github.com/theopenlane/riverboat/test/common"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

// the main function here will insert an email job into the river
// this will be picked up by the river server and processed
// assuming the server is run in the default setup, the email will be written to a file (fixtures/email)
func main() {
	client := common.NewInsertOnlyRiverClient()

	msg := newman.NewEmailMessageWithOptions(
		newman.WithSubject("test subject"),
		newman.WithText("body"),
		newman.WithTo([]string{"meowfunk@example.com"}),
	)

	_, err := client.Insert(context.Background(), jobs.EmailArgs{
		Message: *msg,
	}, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("error inserting email job")
	}

	log.Info().Msg("email job successfully inserted")
}
