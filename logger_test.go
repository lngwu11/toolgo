// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package toolgo_test

import (
	"github.com/lngwu11/toolgo"
	gc "gopkg.in/check.v1"
)

type LoggerSuite struct{}

var _ = gc.Suite(&LoggerSuite{})

func (*LoggerSuite) SetUpTest(c *gc.C) {
	toolgo.ResetDefaultContext()
}

func (s *LoggerSuite) TestRootLogger(c *gc.C) {
	root := toolgo.Logger{}
	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.LogLevel(), gc.Equals, toolgo.WARNING)
	c.Check(root.IsErrorEnabled(), gc.Equals, true)
	c.Check(root.IsWarningEnabled(), gc.Equals, true)
	c.Check(root.IsInfoEnabled(), gc.Equals, false)
	c.Check(root.IsDebugEnabled(), gc.Equals, false)
	c.Check(root.IsTraceEnabled(), gc.Equals, false)
}

func (s *LoggerSuite) TestSetLevel(c *gc.C) {
	logger := toolgo.GetLogger("testing")

	c.Assert(logger.LogLevel(), gc.Equals, toolgo.UNSPECIFIED)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, toolgo.WARNING)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(toolgo.TRACE)
	c.Assert(logger.LogLevel(), gc.Equals, toolgo.TRACE)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, toolgo.TRACE)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, true)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, true)
	logger.SetLogLevel(toolgo.DEBUG)
	c.Assert(logger.LogLevel(), gc.Equals, toolgo.DEBUG)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, toolgo.DEBUG)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, true)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(toolgo.INFO)
	c.Assert(logger.LogLevel(), gc.Equals, toolgo.INFO)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, toolgo.INFO)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(toolgo.WARNING)
	c.Assert(logger.LogLevel(), gc.Equals, toolgo.WARNING)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, toolgo.WARNING)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(toolgo.ERROR)
	c.Assert(logger.LogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, false)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	// This is added for completeness, but not really expected to be used.
	logger.SetLogLevel(toolgo.CRITICAL)
	c.Assert(logger.LogLevel(), gc.Equals, toolgo.CRITICAL)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, toolgo.CRITICAL)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, false)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, false)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(toolgo.UNSPECIFIED)
	c.Assert(logger.LogLevel(), gc.Equals, toolgo.UNSPECIFIED)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, toolgo.WARNING)
}

func (s *LoggerSuite) TestModuleLowered(c *gc.C) {
	logger1 := toolgo.GetLogger("TESTING.MODULE")
	logger2 := toolgo.GetLogger("Testing")

	c.Assert(logger1.Name(), gc.Equals, "testing.module")
	c.Assert(logger2.Name(), gc.Equals, "testing")
}

func (s *LoggerSuite) TestLevelsInherited(c *gc.C) {
	root := toolgo.GetLogger("")
	first := toolgo.GetLogger("first")
	second := toolgo.GetLogger("first.second")

	root.SetLogLevel(toolgo.ERROR)
	c.Assert(root.LogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, toolgo.UNSPECIFIED)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(second.LogLevel(), gc.Equals, toolgo.UNSPECIFIED)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, toolgo.ERROR)

	first.SetLogLevel(toolgo.DEBUG)
	c.Assert(root.LogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, toolgo.DEBUG)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, toolgo.DEBUG)
	c.Assert(second.LogLevel(), gc.Equals, toolgo.UNSPECIFIED)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, toolgo.DEBUG)

	second.SetLogLevel(toolgo.INFO)
	c.Assert(root.LogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, toolgo.DEBUG)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, toolgo.DEBUG)
	c.Assert(second.LogLevel(), gc.Equals, toolgo.INFO)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, toolgo.INFO)

	first.SetLogLevel(toolgo.UNSPECIFIED)
	c.Assert(root.LogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, toolgo.UNSPECIFIED)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, toolgo.ERROR)
	c.Assert(second.LogLevel(), gc.Equals, toolgo.INFO)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, toolgo.INFO)
}

func (s *LoggerSuite) TestParent(c *gc.C) {
	logger := toolgo.GetLogger("a.b.c")
	b := logger.Parent()
	a := b.Parent()
	root := a.Parent()

	c.Check(b.Name(), gc.Equals, "a.b")
	c.Check(a.Name(), gc.Equals, "a")
	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.Parent(), gc.Equals, root)
}

func (s *LoggerSuite) TestParentSameContext(c *gc.C) {
	ctx := toolgo.NewContext(toolgo.DEBUG)

	logger := ctx.GetLogger("a.b.c")
	b := logger.Parent()

	c.Check(b, gc.Equals, ctx.GetLogger("a.b"))
	c.Check(b, gc.Not(gc.Equals), toolgo.GetLogger("a.b"))
}

func (s *LoggerSuite) TestChild(c *gc.C) {
	root := toolgo.GetLogger("")

	a := root.Child("a")
	logger := a.Child("b.c")

	c.Check(a.Name(), gc.Equals, "a")
	c.Check(logger.Name(), gc.Equals, "a.b.c")
	c.Check(logger.Parent(), gc.Equals, a.Child("b"))
}

func (s *LoggerSuite) TestChildSameContext(c *gc.C) {
	ctx := toolgo.NewContext(toolgo.DEBUG)

	logger := ctx.GetLogger("a")
	b := logger.Child("b")

	c.Check(b, gc.Equals, ctx.GetLogger("a.b"))
	c.Check(b, gc.Not(gc.Equals), toolgo.GetLogger("a.b"))
}

func (s *LoggerSuite) TestChildSameContextWithLabels(c *gc.C) {
	ctx := toolgo.NewContext(toolgo.DEBUG)

	logger := ctx.GetLogger("a", "parent")
	b := logger.ChildWithLabels("b", "child")

	c.Check(ctx.GetAllLoggerLabels(), gc.DeepEquals, []string{"child", "parent"})
	c.Check(logger.Labels(), gc.DeepEquals, []string{"parent"})
	c.Check(b.Labels(), gc.DeepEquals, []string{"child"})
}

func (s *LoggerSuite) TestRoot(c *gc.C) {
	logger := toolgo.GetLogger("a.b.c")
	root := logger.Root()

	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.Child("a.b.c"), gc.Equals, logger)
}

func (s *LoggerSuite) TestRootSameContext(c *gc.C) {
	ctx := toolgo.NewContext(toolgo.DEBUG)

	logger := ctx.GetLogger("a.b.c")
	root := logger.Root()

	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.Child("a.b.c"), gc.Equals, logger)
	c.Check(root.Child("a.b.c"), gc.Not(gc.Equals), toolgo.GetLogger("a.b.c"))
}
