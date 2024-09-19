package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/theopenlane/riverboat/test/common"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

// the main function here will insert an database job into the river
// this will be picked up by the river server and processed
func main() {
	client := common.NewInsertOnlyRiverClient()

	_, err := client.Insert(context.Background(), jobs.DatabaseArgs{
		OrganizationID: "100100100010001000",
		Location:       "AMER",
	}, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("error inserting database job")
	}

	log.Info().Msg("database job successfully inserted")
}
