[![Go Report Card](https://goreportcard.com/badge/github.com/theopenlane/riverboat)](https://goreportcard.com/report/github.com/theopenlane/riverboat)
[![Build status](https://badge.buildkite.com/34ad31fe4231b2953cd3f2d116364d21a39b2a4dbf1eea539a.svg)](https://buildkite.com/theopenlane/riverboat?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/theopenlane/riverboat.svg)](https://pkg.go.dev/github.com/theopenlane/riverboat)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache2.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=theopenlane_riverboat&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=theopenlane_riverboat)

# Riverboat

Riverboat is the job queue used in openlane based on the
[riverqueue](https://riverqueue.com/) project.

## Usage

Jobs can be inserted into the job queue either from this server directly, or
from any codebase with an
[Insert Only](https://riverqueue.com/docs#insert-only-clients) river client. All
jobs will be processed via the `riverboat` server. Since jobs are committed to
the database within a transaction, and stored in the database we do not have to
worry about dropped events.

## Getting Started

This repo includes several [Taskfiles](https://taskfile.dev/) to assist with
getting things running.

### Dependencies

- Go 1.25.1+
- Docker (used for running Postgres and the river-ui)
- [task](https://taskfile.dev/)

### Starting the Server

The following will start up postgres, the river-ui, and the riverboat server:

```bash
task run-dev
```

### Test Jobs

Included in the `test/` directory are test jobs corresponding to the job types
in `pkg/jobs`.

1. Start the `riverboat` server using `task run-dev`
1. Run the test main, for example the `email`:

   ```bash
   go run test/email/main.go
   ```

1. This should insert the job successfully, it should be processed by `river`
   and the email should be added to `fixtures/email`

## Adding New Jobs

1. New jobs should be added to the `pkg/jobs` directory in a new file, refer to
   the [upstream docs](https://riverqueue.com/docs#job-args-and-workers) for
   implementation details. The following is a stem job that could be copied to
   get you started.

   ```go
   package jobs

   import (
      "context"

      "github.com/riverqueue/river"
      "github.com/rs/zerolog/log"
   )

   // ExampleArgs for the example worker to process the job
   type ExampleArgs struct {
      // ExampleArg is an example argument
      ExampleArg string `json:"example_arg"`
   }

   // Kind satisfies the river.Job interface
   func (ExampleArgs) Kind() string { return "example" }

   // ExampleWorker does all sorts of neat stuff
   type ExampleWorker struct {
      river.WorkerDefaults[ExampleArgs]

      ExampleConfig
   }

   // ExampleConfig contains the configuration for the example worker
   type ExampleConfig struct {
      // DevMode is a flag to enable dev mode so we don't actually send millions of carrier pigeons
      DevMode bool `koanf:"devMode" json:"devMode" jsonschema:"description=enable dev mode" default:"true"`
   }

   // Work satisfies the river.Worker interface for the example worker
   func (w *ExampleConfig) Work(ctx context.Context, job *river.Job[ExampleArgs]) error {
      // do some work

      return nil
   }
   ```

1. Add a test for the new job, see `email_test.go` as an example. There are
   additional helper functions that can be used, see
   [river test helpers](https://riverqueue.com/docs/testing) for details.
1. If there are configuration settings, add the worker to `pkg/river/config.go`
   `Workers` struct, this will allow the config variables to be set via the
   `koanf` config setup. Once added you will need to regenerate the config:

   ```bash
   task config:generate
   ```

1. Register the worker by adding the `river.AddWorkerSafely` to the
   `pkg/river/workers.go` `createWorkers` function.
1. Add a `test` job to `test/` directory by creating a new directory with a
   `main.go` function that will insert the job into the queue.

## Contributing

See the [contributing](.github/CONTRIBUTING.md) guide for more information.
