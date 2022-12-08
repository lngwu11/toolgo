// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package toolgo_test

import (
	"github.com/lngwu11/toolgo"
	"time"

	gc "gopkg.in/check.v1"
)

type LoggingSuite struct {
	context *toolgo.Context
	writer  *writer
	logger  toolgo.Logger

	// Test that labels get outputted to toolgo.Entry
	Labels []string
}

var _ = gc.Suite(&LoggingSuite{})
var _ = gc.Suite(&LoggingSuite{Labels: []string{"ONE", "TWO"}})

func (s *LoggingSuite) SetUpTest(c *gc.C) {
	s.writer = &writer{}
	s.context = toolgo.NewContext(toolgo.TRACE)
	err := s.context.AddWriter("test", s.writer)
	c.Assert(err, gc.IsNil)
	s.logger = s.context.GetLogger("test", s.Labels...)
}

func (s *LoggingSuite) TestLoggingStrings(c *gc.C) {
	s.logger.Infof("simple")
	s.logger.Infof("with args %d", 42)
	s.logger.Infof("working 100%")
	s.logger.Infof("missing %s")

	checkLogEntries(c, s.writer.Log(), []toolgo.Entry{
		{Level: toolgo.INFO, Module: "test", Message: "simple", Labels: s.Labels},
		{Level: toolgo.INFO, Module: "test", Message: "with args 42", Labels: s.Labels},
		{Level: toolgo.INFO, Module: "test", Message: "working 100%", Labels: s.Labels},
		{Level: toolgo.INFO, Module: "test", Message: "missing %s", Labels: s.Labels},
	})
}

func (s *LoggingSuite) TestLoggingLimitWarning(c *gc.C) {
	s.logger.SetLogLevel(toolgo.WARNING)
	start := time.Now()
	logAllSeverities(s.logger)
	end := time.Now()
	entries := s.writer.Log()
	checkLogEntries(c, entries, []toolgo.Entry{
		{Level: toolgo.CRITICAL, Module: "test", Message: "something critical", Labels: s.Labels},
		{Level: toolgo.ERROR, Module: "test", Message: "an error", Labels: s.Labels},
		{Level: toolgo.WARNING, Module: "test", Message: "a warning message", Labels: s.Labels},
	})

	for _, entry := range entries {
		c.Check(entry.Timestamp, Between(start, end))
	}
}

func (s *LoggingSuite) TestLocationCapture(c *gc.C) {
	s.logger.Criticalf("critical message") //tag critical-location
	s.logger.Errorf("error message")       //tag error-location
	s.logger.Warningf("warning message")   //tag warning-location
	s.logger.Infof("info message")         //tag info-location
	s.logger.Debugf("debug message")       //tag debug-location
	s.logger.Tracef("trace message")       //tag trace-location

	log := s.writer.Log()
	tags := []string{
		"critical-location",
		"error-location",
		"warning-location",
		"info-location",
		"debug-location",
		"trace-location",
	}
	c.Assert(log, gc.HasLen, len(tags))
	for x := range tags {
		assertLocation(c, log[x], tags[x])
	}
}

func (s *LoggingSuite) TestLogDoesntLogWeirdLevels(c *gc.C) {
	s.logger.Logf(toolgo.UNSPECIFIED, "message")
	c.Assert(s.writer.Log(), gc.HasLen, 0)

	s.logger.Logf(toolgo.Level(42), "message")
	c.Assert(s.writer.Log(), gc.HasLen, 0)

	s.logger.Logf(toolgo.CRITICAL+toolgo.Level(1), "message")
	c.Assert(s.writer.Log(), gc.HasLen, 0)
}
