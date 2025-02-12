package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/theopenlane/riverboat/test/common"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

// the main function here will insert an email job into the river
// this will be picked up by the river server and processed
// assuming the server is run in the default setup, the email will be written to a file (fixtures/email)
func main() {
	client := common.NewInsertOnlyRiverClient()

	_, err := client.Insert(context.Background(), jobs.OnboardingArgs{
		Token:          "",                           // update before running test
		OrganizationID: "01JKW55RJ1NDGHHE2DANTQ0DVV", // update before running test
		OwnerID:        "01JKW55Q4EJJHN195S2GFF15VV", // update before running test
		AdminGroupID:   "01JKW55RKGQSXZ627RQJDYRGGF", // update before running test
	}, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("error inserting onboarding job")
	}

	log.Info().Msg("onboarding job successfully inserted")
}
