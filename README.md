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

- Go 1.23
- Docker (used for running Postgres)

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

## Contributing

See the [contributing](.github/CONTRIBUTING.md) guide for more information.
