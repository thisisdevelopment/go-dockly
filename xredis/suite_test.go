package xredis_test

import (
	"testing"

	"github.com/thisisdevelopment/go-dockly/v3/xhelper"
	"github.com/thisisdevelopment/go-dockly/v3/xlogger"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	logger *xlogger.Logger
	helper *xhelper.Helper
}

func (s *TestSuite) SetupSuite() {
	s.logger = xlogger.DefaultTestLogger(&s.Suite)
	s.helper = xhelper.NewHelper(&s.Suite, s.logger)
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
