// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package toolgo_test

import (
	"github.com/lngwu11/toolgo"
	"io/ioutil"
	"os"

	gc "gopkg.in/check.v1"
)

type BenchmarksSuite struct {
	logger toolgo.Logger
	writer *writer
}

var _ = gc.Suite(&BenchmarksSuite{})

func (s *BenchmarksSuite) SetUpTest(c *gc.C) {
	toolgo.ResetLogging()
	s.logger = toolgo.GetLogger("test.writer")
	s.writer = &writer{}
	err := toolgo.RegisterWriter("test", s.writer)
	c.Assert(err, gc.IsNil)
}

func (s *BenchmarksSuite) BenchmarkLoggingNoWriters(c *gc.C) {
	// No writers
	_, _ = toolgo.RemoveWriter("test")
	for i := 0; i < c.N; i++ {
		s.logger.Warningf("just a simple warning for %d", i)
	}
}

func (s *BenchmarksSuite) BenchmarkLoggingNoWritersNoFormat(c *gc.C) {
	// No writers
	_, _ = toolgo.RemoveWriter("test")
	for i := 0; i < c.N; i++ {
		s.logger.Warningf("just a simple warning")
	}
}

func (s *BenchmarksSuite) BenchmarkLoggingTestWriters(c *gc.C) {
	for i := 0; i < c.N; i++ {
		s.logger.Warningf("just a simple warning for %d", i)
	}
	c.Assert(s.writer.Log(), gc.HasLen, c.N)
}

func (s *BenchmarksSuite) BenchmarkLoggingDiskWriter(c *gc.C) {
	logFile := s.setupTempFileWriter(c)
	defer logFile.Close()
	msg := "just a simple warning for %d"
	for i := 0; i < c.N; i++ {
		s.logger.Warningf(msg, i)
	}
	offset, err := logFile.Seek(0, os.SEEK_CUR)
	c.Assert(err, gc.IsNil)
	c.Assert((offset > int64(len(msg))*int64(c.N)), gc.Equals, true,
		gc.Commentf("Not enough data was written to the log file."))
}

func (s *BenchmarksSuite) BenchmarkLoggingDiskWriterNoMessages(c *gc.C) {
	logFile := s.setupTempFileWriter(c)
	defer logFile.Close()
	// Change the log level
	writer, err := toolgo.RemoveWriter("testfile")
	c.Assert(err, gc.IsNil)

	err = toolgo.RegisterWriter("testfile", toolgo.NewMinimumLevelWriter(writer, toolgo.WARNING))
	c.Assert(err, gc.IsNil)

	msg := "just a simple warning for %d"
	for i := 0; i < c.N; i++ {
		s.logger.Debugf(msg, i)
	}
	offset, err := logFile.Seek(0, os.SEEK_CUR)
	c.Assert(err, gc.IsNil)
	c.Assert(offset, gc.Equals, int64(0),
		gc.Commentf("Data was written to the log file."))
}

func (s *BenchmarksSuite) BenchmarkLoggingDiskWriterNoMessagesLogLevel(c *gc.C) {
	logFile := s.setupTempFileWriter(c)
	defer logFile.Close()
	// Change the log level
	s.logger.SetLogLevel(toolgo.WARNING)

	msg := "just a simple warning for %d"
	for i := 0; i < c.N; i++ {
		s.logger.Debugf(msg, i)
	}
	offset, err := logFile.Seek(0, os.SEEK_CUR)
	c.Assert(err, gc.IsNil)
	c.Assert(offset, gc.Equals, int64(0),
		gc.Commentf("Data was written to the log file."))
}

func (s *BenchmarksSuite) setupTempFileWriter(c *gc.C) *os.File {
	_, _ = toolgo.RemoveWriter("test")
	logFile, err := ioutil.TempFile(c.MkDir(), "toolgo-test")
	c.Assert(err, gc.IsNil)

	writer := toolgo.NewSimpleWriter(logFile, toolgo.DefaultFormatter)
	err = toolgo.RegisterWriter("testfile", writer)
	c.Assert(err, gc.IsNil)
	return logFile
}
