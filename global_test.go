// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package toolgo_test

import (
	"github.com/lngwu11/toolgo"
	gc "gopkg.in/check.v1"
)

type GlobalSuite struct{}

var _ = gc.Suite(&GlobalSuite{})

func (*GlobalSuite) SetUpTest(c *gc.C) {
	toolgo.ResetDefaultContext()
}

func (*GlobalSuite) TestRootLogger(c *gc.C) {
	var root toolgo.Logger

	got := toolgo.GetLogger("")

	c.Check(got.Name(), gc.Equals, root.Name())
	c.Check(got.LogLevel(), gc.Equals, root.LogLevel())
}

func (*GlobalSuite) TestModuleName(c *gc.C) {
	logger := toolgo.GetLogger("toolgo.testing")
	c.Check(logger.Name(), gc.Equals, "toolgo.testing")
}

func (*GlobalSuite) TestLevel(c *gc.C) {
	logger := toolgo.GetLogger("testing")
	level := logger.LogLevel()
	c.Check(level, gc.Equals, toolgo.UNSPECIFIED)
}

func (*GlobalSuite) TestEffectiveLevel(c *gc.C) {
	logger := toolgo.GetLogger("testing")
	level := logger.EffectiveLogLevel()
	c.Check(level, gc.Equals, toolgo.WARNING)
}

func (*GlobalSuite) TestLevelsSharedForSameModule(c *gc.C) {
	logger1 := toolgo.GetLogger("testing.module")
	logger2 := toolgo.GetLogger("testing.module")

	logger1.SetLogLevel(toolgo.INFO)
	c.Assert(logger1.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger2.IsInfoEnabled(), gc.Equals, true)
}

func (*GlobalSuite) TestModuleLowered(c *gc.C) {
	logger1 := toolgo.GetLogger("TESTING.MODULE")
	logger2 := toolgo.GetLogger("Testing")

	c.Assert(logger1.Name(), gc.Equals, "testing.module")
	c.Assert(logger2.Name(), gc.Equals, "testing")
}

func (s *GlobalSuite) TestConfigureLoggers(c *gc.C) {
	err := toolgo.ConfigureLoggers("testing.module=debug")
	c.Assert(err, gc.IsNil)
	expected := "<root>=WARNING;testing.module=DEBUG"
	c.Assert(toolgo.DefaultContext().Config().String(), gc.Equals, expected)
	c.Assert(toolgo.LoggerInfo(), gc.Equals, expected)
}

func (*GlobalSuite) TestRegisterWriterExistingName(c *gc.C) {
	err := toolgo.RegisterWriter("default", &writer{})
	c.Assert(err, gc.ErrorMatches, `context already has a writer named "default"`)
}

func (*GlobalSuite) TestReplaceDefaultWriter(c *gc.C) {
	oldWriter, err := toolgo.ReplaceDefaultWriter(&writer{})
	c.Assert(oldWriter, gc.NotNil)
	c.Assert(err, gc.IsNil)
	c.Assert(toolgo.DefaultContext().WriterNames(), gc.DeepEquals, []string{"default"})
}

func (*GlobalSuite) TestRemoveWriter(c *gc.C) {
	oldWriter, err := toolgo.RemoveWriter("default")
	c.Assert(oldWriter, gc.NotNil)
	c.Assert(err, gc.IsNil)
	c.Assert(toolgo.DefaultContext().WriterNames(), gc.HasLen, 0)
}

func (s *GlobalSuite) TestGetLoggerWithLabels(c *gc.C) {
	logger := toolgo.GetLoggerWithLabels("parent", "labela", "labelb")
	c.Check(logger.Labels(), gc.DeepEquals, []string{"labela", "labelb"})
}
