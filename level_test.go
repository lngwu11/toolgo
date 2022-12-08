// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package toolgo_test

import (
	"github.com/lngwu11/toolgo"
	gc "gopkg.in/check.v1"
)

type LevelSuite struct{}

var _ = gc.Suite(&LevelSuite{})

var parseLevelTests = []struct {
	str   string
	level toolgo.Level
	fail  bool
}{{
	str:   "trace",
	level: toolgo.TRACE,
}, {
	str:   "TrAce",
	level: toolgo.TRACE,
}, {
	str:   "TRACE",
	level: toolgo.TRACE,
}, {
	str:   "debug",
	level: toolgo.DEBUG,
}, {
	str:   "DEBUG",
	level: toolgo.DEBUG,
}, {
	str:   "info",
	level: toolgo.INFO,
}, {
	str:   "INFO",
	level: toolgo.INFO,
}, {
	str:   "warn",
	level: toolgo.WARNING,
}, {
	str:   "WARN",
	level: toolgo.WARNING,
}, {
	str:   "warning",
	level: toolgo.WARNING,
}, {
	str:   "WARNING",
	level: toolgo.WARNING,
}, {
	str:   "error",
	level: toolgo.ERROR,
}, {
	str:   "ERROR",
	level: toolgo.ERROR,
}, {
	str:   "critical",
	level: toolgo.CRITICAL,
}, {
	str:  "not_specified",
	fail: true,
}, {
	str:  "other",
	fail: true,
}, {
	str:  "",
	fail: true,
}}

func (s *LevelSuite) TestParseLevel(c *gc.C) {
	for _, test := range parseLevelTests {
		level, ok := toolgo.ParseLevel(test.str)
		c.Assert(level, gc.Equals, test.level)
		c.Assert(ok, gc.Equals, !test.fail)
	}
}

var levelStringValueTests = map[toolgo.Level]string{
	toolgo.UNSPECIFIED: "UNSPECIFIED",
	toolgo.DEBUG:       "DEBUG",
	toolgo.TRACE:       "TRACE",
	toolgo.INFO:        "INFO",
	toolgo.WARNING:     "WARNING",
	toolgo.ERROR:       "ERROR",
	toolgo.CRITICAL:    "CRITICAL",
	toolgo.Level(42):   "<unknown>", // other values are unknown
}

func (s *LevelSuite) TestLevelStringValue(c *gc.C) {
	for level, str := range levelStringValueTests {
		c.Assert(level.String(), gc.Equals, str)
	}
}
