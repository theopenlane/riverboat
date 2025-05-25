package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/theopenlane/riverboat/pkg/jobs"
	"github.com/theopenlane/riverboat/test/common"
)

// the main function here will insert a create_custom_domain job into the river
// this will be picked up by the river server and processed
func main() {
	client := common.NewInsertOnlyRiverClient()

	_, err := client.Insert(context.Background(), jobs.CreateCustomDomainArgs{
		CustomDomainID: "test-domain-123",
	}, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("error inserting create_custom_domain job")
	}

	log.Info().Msg("create_custom_domain job successfully inserted")
}
