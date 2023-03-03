package loggo

import "time"

// Entry represents a single log message.
type Entry struct {
	// Level is the severity of the log message.
	Level Level
	// Module is the dotted module name from the logger.
	Module string
	// Filename is the full path the file that logged the message.
	Filename string
	// Line is the line number of the Filename.
	Line int
	// Timestamp is when the log message was created
	Timestamp time.Time
	// Message is the formatted string from teh log call.
	Message string
	// Labels is the label associated with the log message.
	Labels []string
}
