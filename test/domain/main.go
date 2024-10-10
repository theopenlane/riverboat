package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/theopenlane/riverboat/test/common"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

// the main function here will insert an domain job into the river
// this will be picked up by the river server and processed
// assuming the server is run in the default setup, the domain worker will check the domain for subdomains
func main() {
	client := common.NewInsertOnlyRiverClient()

	domain := "theopenlane.io"

	_, err := client.Insert(context.Background(), jobs.DomainArgs{
		Domain: domain,
	}, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("error inserting domain job")
	}

	log.Info().Msg("domain job successfully inserted")
}
