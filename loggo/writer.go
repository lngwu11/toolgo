package loggo

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// DefaultWriterName is the name of the default writer for
// a Context.
const DefaultWriterName = "default"

// Writer is implemented by any recipient of log messages.
type Writer interface {
	// Write writes a message to the Writer with the given level and module
	// name. The filename and line hold the file name and line number of the
	// code that is generating the log message; the time stamp holds the time
	// the log message was generated, and message holds the log message
	// itself.
	Write(entry Entry)
}

// NewMinimumLevelWriter returns a Writer that will only pass on the Write calls
// to the provided writer if the log level is at or above the specified
// minimum level.
func NewMinimumLevelWriter(writer Writer, minLevel Level) Writer {
	return &minLevelWriter{
		writer: writer,
		level:  minLevel,
	}
}

type minLevelWriter struct {
	writer Writer
	level  Level
}

// Write writes the log record.
func (w minLevelWriter) Write(entry Entry) {
	if entry.Level < w.level {
		return
	}
	w.writer.Write(entry)
}

type simpleWriter struct {
	writer    io.Writer
	formatter func(entry Entry) string
}

// NewSimpleWriter returns a new writer that writes log messages to the given
// io.Writer formatting the messages with the given formatter.
func NewSimpleWriter(writer io.Writer, formatter func(entry Entry) string) Writer {
	if formatter == nil {
		formatter = DefaultFormatter
	}
	return &simpleWriter{writer, formatter}
}

func (simple *simpleWriter) Write(entry Entry) {
	logLine := simple.formatter(entry)
	_, _ = fmt.Fprintln(simple.writer, logLine)
}

func defaultWriter() Writer {
	return NewSimpleWriter(os.Stderr, DefaultFormatter)
}

// LogFileWriter is a simple log file writer
type LogFileWriter struct {
	name      string
	perm      os.FileMode
	formatter func(entry Entry) string
	reopen    chan os.Signal
	file      *os.File
}

func (logfile *LogFileWriter) Write(entry Entry) {
	select {
	case _, ok := <-logfile.reopen:
		if ok {
			if file, err := os.OpenFile(logfile.name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, logfile.perm); err == nil {
				_ = logfile.file.Close()
				logfile.file = file
			}
		} else {
			panic("logfile already closed")
		}
	default:
	}

	logLine := logfile.formatter(entry)
	_, _ = fmt.Fprintln(logfile.file, logLine)
}

// Close a log file writer
func (logfile *LogFileWriter) Close() error {
	signal.Stop(logfile.reopen)
	close(logfile.reopen)
	return logfile.file.Close()
}

// NewLogFileWriter returns a new writer that writes log messages to the given
// file and formatting the messages with the given formatter.
func NewLogFileWriter(name string, perm os.FileMode, formatter func(entry Entry) string) (*LogFileWriter, error) {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return nil, err
	}
	if formatter == nil {
		formatter = DefaultFormatter
	}
	writer := &LogFileWriter{
		name:      name,
		perm:      perm,
		formatter: formatter,
		reopen:    make(chan os.Signal, 1),
		file:      file,
	}
	signal.Notify(writer.reopen, syscall.SIGHUP)
	return writer, nil
}
