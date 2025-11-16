package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/riverboat/pkg/jobs"
	"github.com/theopenlane/riverboat/test/common"
)

func main() {
	client := common.NewInsertOnlyRiverClient()
	args := jobs.SlackArgs{
		Channel: "general", // or channel ID
		Message: "Hello from riverboat test job!",
		DevMode: true, // set to true in dev mode, change in prod
	}
	_, err := client.Insert(context.Background(), args, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("error inserting slack job")
	}

	log.Info().Msg("slack job successfully inserted")
}
