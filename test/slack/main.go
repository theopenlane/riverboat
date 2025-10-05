package main

import (
	"context"
	"log"

	"github.com/riverqueue/river"
	"github.com/theopenlane/riverboat/pkg/jobs"
)

func main() {
	// Inserting a Slack job into the queue
	args := jobs.SlackArgs{
		Channel: "general", // or channel ID
		Message: "Hello from riverboat test job!",
		DevMode: true, // set to true in dev mode, please change it in prod.
	}
	job := river.NewJob(args)
	if err := river.Insert(context.Background(), job); err != nil {
		log.Fatalf("failed to insert slack job: %v", err)
	}
	log.Println("Inserted slack job!")
}
