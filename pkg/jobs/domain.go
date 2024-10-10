package jobs

import (
	"bytes"
	"context"
	"io"
	"strings"

	subfinder "github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
)

type DomainArgs struct {
	// Domain is the domain to check
	Domain string `json:"domain"`
}

// Kind satisfies the river.Job interface
func (DomainArgs) Kind() string { return "domain" }

// DomainWorker is a worker to check the domain using the domain provider
type DomainWorker struct {
	river.WorkerDefaults[DomainArgs]

	Config SubFinderConfig `koanf:"config" json:"config" jsonschema:"description=the subfinder configuration"`
}

type SubFinderConfig struct {
	// Threads is the number of threads to use
	Threads int `koanf:"threads" json:"threads" jsonschema:"description=the number of threads to use" default:"10"`
	// Timeout is the timeout to use
	Timeout int `koanf:"timeout" json:"timeout" jsonschema:"description=the timeout to use" default:"30"`
	// MaxEnumerationTime is the max enumeration time to use
	MaxEnumerationTime int `koanf:"maxEnumerationTime" json:"maxEnumerationTime" jsonschema:"description=the max enumeration time to use" default:"5"`
}

// validateDomainInput validates the input for the domain worker
func validateDomainInput(job *river.Job[DomainArgs]) error {
	if job.Args.Domain == "" {
		return newMissingRequiredArg("domain", DomainArgs{}.Kind())
	}

	return nil
}

// Work satisfies the river.Worker interface for the domain worker
// it checks the domain for all subdomains
func (w *DomainWorker) Work(ctx context.Context, job *river.Job[DomainArgs]) error {
	if err := validateDomainInput(job); err != nil {
		return err
	}

	opts := subfinder.Options{
		Threads:            w.Config.Threads,
		Timeout:            w.Config.Timeout,
		MaxEnumerationTime: w.Config.MaxEnumerationTime,
	}

	sf, err := subfinder.NewRunner(&opts)
	if err != nil {
		return err
	}

	out := &bytes.Buffer{}

	// check the domain for all subdomains
	if err := sf.EnumerateSingleDomainWithCtx(ctx, job.Args.Domain, []io.Writer{out}); err != nil {
		return err
	}

	f := func(c rune) bool {
		return c == '\n'
	}

	subdomains := strings.FieldsFunc(out.String(), f)

	log.Info().Str("domain", job.Args.Domain).Strs("subdomains", subdomains).Msg("subdomains found")

	return nil
}
