// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package toolgo_test

import (
	"github.com/lngwu11/toolgo"
	"time"

	gc "gopkg.in/check.v1"
)

type formatterSuite struct{}

var _ = gc.Suite(&formatterSuite{})

func (*formatterSuite) TestDefaultFormat(c *gc.C) {
	location, err := time.LoadLocation("UTC")
	testTime := time.Date(2013, 5, 3, 10, 53, 24, 123456, location)
	c.Assert(err, gc.IsNil)
	entry := toolgo.Entry{
		Level:     toolgo.WARNING,
		Module:    "test.module",
		Filename:  "some/deep/filename",
		Line:      42,
		Timestamp: testTime,
		Message:   "hello world!",
	}
	formatted := toolgo.DefaultFormatter(entry)
	c.Assert(formatted, gc.Equals, "2013-05-03 10:53:24 WARNING test.module filename:42 hello world!")
}
