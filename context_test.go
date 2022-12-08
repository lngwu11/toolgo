// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package toolgo_test

import (
	"github.com/lngwu11/toolgo"
	gc "gopkg.in/check.v1"
)

type ContextSuite struct{}

var _ = gc.Suite(&ContextSuite{})

func (*ContextSuite) TestNewContextRootLevel(c *gc.C) {
	for i, test := range []struct {
		level    toolgo.Level
		expected toolgo.Level
	}{{
		level:    toolgo.UNSPECIFIED,
		expected: toolgo.WARNING,
	}, {
		level:    toolgo.DEBUG,
		expected: toolgo.DEBUG,
	}, {
		level:    toolgo.INFO,
		expected: toolgo.INFO,
	}, {
		level:    toolgo.WARNING,
		expected: toolgo.WARNING,
	}, {
		level:    toolgo.ERROR,
		expected: toolgo.ERROR,
	}, {
		level:    toolgo.CRITICAL,
		expected: toolgo.CRITICAL,
	}, {
		level:    toolgo.Level(42),
		expected: toolgo.WARNING,
	}} {
		c.Logf("%d: %s", i, test.level)
		context := toolgo.NewContext(test.level)
		cfg := context.Config()
		c.Check(cfg, gc.HasLen, 1)
		value, found := cfg[""]
		c.Check(found, gc.Equals, true)
		c.Check(value, gc.Equals, test.expected)
	}
}

func logAllSeverities(logger toolgo.Logger) {
	logger.Criticalf("something critical")
	logger.Errorf("an error")
	logger.Warningf("a warning message")
	logger.Infof("an info message")
	logger.Debugf("a debug message")
	logger.Tracef("a trace message")
}

func checkLogEntry(c *gc.C, entry, expected toolgo.Entry) {
	c.Check(entry.Level, gc.Equals, expected.Level)
	c.Check(entry.Module, gc.Equals, expected.Module)
	c.Check(entry.Message, gc.Equals, expected.Message)
}

func checkLogEntries(c *gc.C, obtained, expected []toolgo.Entry) {
	if c.Check(len(obtained), gc.Equals, len(expected)) {
		for i := range obtained {
			checkLogEntry(c, obtained[i], expected[i])
		}
	}
}

func (*ContextSuite) TestGetLoggerRoot(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	blank := context.GetLogger("")
	root := context.GetLogger("<root>")
	c.Assert(blank, gc.Equals, root)
}

func (*ContextSuite) TestGetLoggerCase(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	upper := context.GetLogger("TEST")
	lower := context.GetLogger("test")
	c.Assert(upper, gc.Equals, lower)
	c.Assert(upper.Name(), gc.Equals, "test")
}

func (*ContextSuite) TestGetLoggerSpace(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	space := context.GetLogger(" test ")
	lower := context.GetLogger("test")
	c.Assert(space, gc.Equals, lower)
	c.Assert(space.Name(), gc.Equals, "test")
}

func (*ContextSuite) TestNewContextNoWriter(c *gc.C) {
	// Should be no output.
	context := toolgo.NewContext(toolgo.DEBUG)
	logger := context.GetLogger("test")
	logAllSeverities(logger)
}

func (*ContextSuite) newContextWithTestWriter(c *gc.C, level toolgo.Level) (*toolgo.Context, *toolgo.TestWriter) {
	writer := &toolgo.TestWriter{}
	context := toolgo.NewContext(level)
	err := context.AddWriter("test", writer)
	c.Assert(err, gc.IsNil)
	return context, writer
}

func (s *ContextSuite) TestNewContextRootSeverityWarning(c *gc.C) {
	context, writer := s.newContextWithTestWriter(c, toolgo.WARNING)
	logger := context.GetLogger("test")
	logAllSeverities(logger)
	checkLogEntries(c, writer.Log(), []toolgo.Entry{
		{Level: toolgo.CRITICAL, Module: "test", Message: "something critical"},
		{Level: toolgo.ERROR, Module: "test", Message: "an error"},
		{Level: toolgo.WARNING, Module: "test", Message: "a warning message"},
	})
}

func (s *ContextSuite) TestNewContextRootSeverityTrace(c *gc.C) {
	context, writer := s.newContextWithTestWriter(c, toolgo.TRACE)
	logger := context.GetLogger("test")
	logAllSeverities(logger)
	checkLogEntries(c, writer.Log(), []toolgo.Entry{
		{Level: toolgo.CRITICAL, Module: "test", Message: "something critical"},
		{Level: toolgo.ERROR, Module: "test", Message: "an error"},
		{Level: toolgo.WARNING, Module: "test", Message: "a warning message"},
		{Level: toolgo.INFO, Module: "test", Message: "an info message"},
		{Level: toolgo.DEBUG, Module: "test", Message: "a debug message"},
		{Level: toolgo.TRACE, Module: "test", Message: "a trace message"},
	})
}

func (*ContextSuite) TestNewContextConfig(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	config := context.Config()
	c.Assert(config, gc.DeepEquals, toolgo.Config{"": toolgo.DEBUG})
}

func (*ContextSuite) TestNewLoggerAddsConfig(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	_ = context.GetLogger("test.module")
	c.Assert(context.Config(), gc.DeepEquals, toolgo.Config{
		"": toolgo.DEBUG,
	})
	c.Assert(context.CompleteConfig(), gc.DeepEquals, toolgo.Config{
		"":            toolgo.DEBUG,
		"test":        toolgo.UNSPECIFIED,
		"test.module": toolgo.UNSPECIFIED,
	})
}

func (*ContextSuite) TestConfigureLoggers(c *gc.C) {
	context := toolgo.NewContext(toolgo.INFO)
	err := context.ConfigureLoggers("testing.module=debug")
	c.Assert(err, gc.IsNil)
	expected := "<root>=INFO;testing.module=DEBUG"
	c.Assert(context.Config().String(), gc.Equals, expected)
}

func (*ContextSuite) TestApplyNilConfig(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	context.ApplyConfig(nil)
	c.Assert(context.Config(), gc.DeepEquals, toolgo.Config{"": toolgo.DEBUG})
}

func (*ContextSuite) TestApplyConfigRootUnspecified(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	context.ApplyConfig(toolgo.Config{"": toolgo.UNSPECIFIED})
	c.Assert(context.Config(), gc.DeepEquals, toolgo.Config{"": toolgo.WARNING})
}

func (*ContextSuite) TestApplyConfigRootTrace(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)
	context.ApplyConfig(toolgo.Config{"": toolgo.TRACE})
	c.Assert(context.Config(), gc.DeepEquals, toolgo.Config{"": toolgo.TRACE})
}

func (*ContextSuite) TestApplyConfigCreatesModules(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)
	context.ApplyConfig(toolgo.Config{"first.second": toolgo.TRACE})
	c.Assert(context.Config(), gc.DeepEquals,
		toolgo.Config{
			"":             toolgo.WARNING,
			"first.second": toolgo.TRACE,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		toolgo.Config{
			"":             toolgo.WARNING,
			"first":        toolgo.UNSPECIFIED,
			"first.second": toolgo.TRACE,
		})
}

func (*ContextSuite) TestApplyConfigAdditive(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)
	context.ApplyConfig(toolgo.Config{"first.second": toolgo.TRACE})
	context.ApplyConfig(toolgo.Config{"other.module": toolgo.DEBUG})
	c.Assert(context.Config(), gc.DeepEquals,
		toolgo.Config{
			"":             toolgo.WARNING,
			"first.second": toolgo.TRACE,
			"other.module": toolgo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		toolgo.Config{
			"":             toolgo.WARNING,
			"first":        toolgo.UNSPECIFIED,
			"first.second": toolgo.TRACE,
			"other":        toolgo.UNSPECIFIED,
			"other.module": toolgo.DEBUG,
		})
}

func (*ContextSuite) TestGetAllLoggerLabels(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)
	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	labels := context.GetAllLoggerLabels()
	c.Assert(labels, gc.DeepEquals, []string{"one", "two"})
}

func (*ContextSuite) TestGetAllLoggerLabelsWithApplyConfig(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)
	context.ApplyConfig(toolgo.Config{"#one": toolgo.TRACE})

	labels := context.GetAllLoggerLabels()
	c.Assert(labels, gc.DeepEquals, []string{})
}

func (*ContextSuite) TestApplyConfigLabels(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)
	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	context.ApplyConfig(toolgo.Config{"#one": toolgo.TRACE})
	context.ApplyConfig(toolgo.Config{"#two": toolgo.DEBUG})

	c.Assert(context.Config(), gc.DeepEquals,
		toolgo.Config{
			"":    toolgo.WARNING,
			"a.b": toolgo.TRACE,
			"c.d": toolgo.TRACE,
			"e":   toolgo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		toolgo.Config{
			"":    toolgo.WARNING,
			"a":   toolgo.UNSPECIFIED,
			"a.b": toolgo.TRACE,
			"c":   toolgo.UNSPECIFIED,
			"c.d": toolgo.TRACE,
			"e":   toolgo.DEBUG,
		})
}

func (*ContextSuite) TestApplyConfigLabelsAppliesToNewLoggers(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)

	context.ApplyConfig(toolgo.Config{"#one": toolgo.TRACE})
	context.ApplyConfig(toolgo.Config{"#two": toolgo.DEBUG})

	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	c.Assert(context.Config(), gc.DeepEquals,
		toolgo.Config{
			"":    toolgo.WARNING,
			"a.b": toolgo.TRACE,
			"c.d": toolgo.TRACE,
			"e":   toolgo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		toolgo.Config{
			"":    toolgo.WARNING,
			"a":   toolgo.UNSPECIFIED,
			"a.b": toolgo.TRACE,
			"c":   toolgo.UNSPECIFIED,
			"c.d": toolgo.TRACE,
			"e":   toolgo.DEBUG,
		})
}

func (*ContextSuite) TestApplyConfigLabelsAppliesToNewLoggersWithMultipleTags(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)

	// Invert the order here, to ensure that the config order doesn't matter,
	// but the way the tags are ordered in `GetLogger`.
	context.ApplyConfig(toolgo.Config{"#two": toolgo.DEBUG})
	context.ApplyConfig(toolgo.Config{"#one": toolgo.TRACE})

	context.GetLogger("a.b", "one", "two")

	c.Assert(context.Config(), gc.DeepEquals,
		toolgo.Config{
			"":    toolgo.WARNING,
			"a.b": toolgo.TRACE,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		toolgo.Config{
			"":    toolgo.WARNING,
			"a":   toolgo.UNSPECIFIED,
			"a.b": toolgo.TRACE,
		})
}

func (*ContextSuite) TestApplyConfigLabelsResetLoggerLevels(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)

	context.ApplyConfig(toolgo.Config{"#one": toolgo.TRACE})
	context.ApplyConfig(toolgo.Config{"#two": toolgo.DEBUG})

	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	context.ResetLoggerLevels()

	c.Assert(context.Config(), gc.DeepEquals,
		toolgo.Config{
			"": toolgo.WARNING,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		toolgo.Config{
			"":    toolgo.WARNING,
			"a":   toolgo.UNSPECIFIED,
			"a.b": toolgo.UNSPECIFIED,
			"c":   toolgo.UNSPECIFIED,
			"c.d": toolgo.UNSPECIFIED,
			"e":   toolgo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestApplyConfigLabelsAddative(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)
	context.ApplyConfig(toolgo.Config{"#one": toolgo.TRACE})
	context.ApplyConfig(toolgo.Config{"#two": toolgo.DEBUG})
	c.Assert(context.Config(), gc.DeepEquals,
		toolgo.Config{
			"": toolgo.WARNING,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		toolgo.Config{
			"": toolgo.WARNING,
		})
}

func (*ContextSuite) TestApplyConfigWithMalformedLabel(c *gc.C) {
	context := toolgo.NewContext(toolgo.WARNING)
	context.GetLogger("a.b", "one")

	context.ApplyConfig(toolgo.Config{"#ONE.1": toolgo.TRACE})

	c.Assert(context.Config(), gc.DeepEquals,
		toolgo.Config{
			"": toolgo.WARNING,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		toolgo.Config{
			"":    toolgo.WARNING,
			"a":   toolgo.UNSPECIFIED,
			"a.b": toolgo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestResetLoggerLevels(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	context.ApplyConfig(toolgo.Config{"first.second": toolgo.TRACE})
	context.ResetLoggerLevels()
	c.Assert(context.Config(), gc.DeepEquals,
		toolgo.Config{
			"": toolgo.WARNING,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		toolgo.Config{
			"":             toolgo.WARNING,
			"first":        toolgo.UNSPECIFIED,
			"first.second": toolgo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestWriterNamesNone(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	writers := context.WriterNames()
	c.Assert(writers, gc.HasLen, 0)
}

func (*ContextSuite) TestAddWriterNoName(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	err := context.AddWriter("", nil)
	c.Assert(err.Error(), gc.Equals, "name cannot be empty")
}

func (*ContextSuite) TestAddWriterNil(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	err := context.AddWriter("foo", nil)
	c.Assert(err.Error(), gc.Equals, "writer cannot be nil")
}

func (*ContextSuite) TestNamedAddWriter(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	err := context.AddWriter("foo", &writer{name: "foo"})
	c.Assert(err, gc.IsNil)
	err = context.AddWriter("foo", &writer{name: "foo"})
	c.Assert(err.Error(), gc.Equals, `context already has a writer named "foo"`)

	writers := context.WriterNames()
	c.Assert(writers, gc.DeepEquals, []string{"foo"})
}

func (*ContextSuite) TestRemoveWriter(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	w, err := context.RemoveWriter("unknown")
	c.Assert(err.Error(), gc.Equals, `context has no writer named "unknown"`)
	c.Assert(w, gc.IsNil)
}

func (*ContextSuite) TestRemoveWriterFound(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	original := &writer{name: "foo"}
	err := context.AddWriter("foo", original)
	c.Assert(err, gc.IsNil)
	existing, err := context.RemoveWriter("foo")
	c.Assert(err, gc.IsNil)
	c.Assert(existing, gc.Equals, original)

	writers := context.WriterNames()
	c.Assert(writers, gc.HasLen, 0)
}

func (*ContextSuite) TestReplaceWriterNoName(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	existing, err := context.ReplaceWriter("", nil)
	c.Assert(err.Error(), gc.Equals, "name cannot be empty")
	c.Assert(existing, gc.IsNil)
}

func (*ContextSuite) TestReplaceWriterNil(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	existing, err := context.ReplaceWriter("foo", nil)
	c.Assert(err.Error(), gc.Equals, "writer cannot be nil")
	c.Assert(existing, gc.IsNil)
}

func (*ContextSuite) TestReplaceWriterNotFound(c *gc.C) {
	context := toolgo.NewContext(toolgo.DEBUG)
	existing, err := context.ReplaceWriter("foo", &writer{})
	c.Assert(err.Error(), gc.Equals, `context has no writer named "foo"`)
	c.Assert(existing, gc.IsNil)
}

func (*ContextSuite) TestMultipleWriters(c *gc.C) {
	first := &writer{}
	second := &writer{}
	third := &writer{}
	context := toolgo.NewContext(toolgo.TRACE)
	err := context.AddWriter("first", first)
	c.Assert(err, gc.IsNil)
	err = context.AddWriter("second", second)
	c.Assert(err, gc.IsNil)
	err = context.AddWriter("third", third)
	c.Assert(err, gc.IsNil)

	logger := context.GetLogger("test")
	logAllSeverities(logger)

	expected := []toolgo.Entry{
		{Level: toolgo.CRITICAL, Module: "test", Message: "something critical"},
		{Level: toolgo.ERROR, Module: "test", Message: "an error"},
		{Level: toolgo.WARNING, Module: "test", Message: "a warning message"},
		{Level: toolgo.INFO, Module: "test", Message: "an info message"},
		{Level: toolgo.DEBUG, Module: "test", Message: "a debug message"},
		{Level: toolgo.TRACE, Module: "test", Message: "a trace message"},
	}

	checkLogEntries(c, first.Log(), expected)
	checkLogEntries(c, second.Log(), expected)
	checkLogEntries(c, third.Log(), expected)
}

func (*ContextSuite) TestWriter(c *gc.C) {
	first := &writer{name: "first"}
	second := &writer{name: "second"}
	context := toolgo.NewContext(toolgo.TRACE)
	err := context.AddWriter("first", first)
	c.Assert(err, gc.IsNil)
	err = context.AddWriter("second", second)
	c.Assert(err, gc.IsNil)

	c.Check(context.Writer("unknown"), gc.IsNil)
	c.Check(context.Writer("first"), gc.Equals, first)
	c.Check(context.Writer("second"), gc.Equals, second)

	c.Check(first, gc.Not(gc.Equals), second)
}

type writer struct {
	toolgo.TestWriter
	// The name exists to discriminate writer equality.
	name string
}
