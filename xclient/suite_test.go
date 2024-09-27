package xclient_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/thisisdevelopment/go-dockly/v3/xclient"
	"github.com/thisisdevelopment/go-dockly/v3/xhelper"
	"github.com/thisisdevelopment/go-dockly/v3/xlogger"
)

type TestSuite struct {
	suite.Suite
	logger  *xlogger.Logger
	helper  *xhelper.Helper
	client  xclient.IAPIClient
	baseURL string
}

func (s *TestSuite) SetupSuite() {
	s.helper = xhelper.NewHelper(&s.Suite, s.logger)

	s.logger = s.helper.GetLogger()
	s.baseURL = "http://test.com"

	cli, err := xclient.New(s.logger, s.baseURL, nil, nil)
	require.NoError(s.T(), err)

	s.client = cli
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
