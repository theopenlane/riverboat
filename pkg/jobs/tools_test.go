package jobs_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

// TestGraphTestSuite runs all the tests in the GraphTestSuite
func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// TestSuite handles the setup and teardown between tests
type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func (suite *TestSuite) TearDownSuite() {
}

func (suite *TestSuite) SetupTest() {
}

func (suite *TestSuite) TearDownTest() {
}
