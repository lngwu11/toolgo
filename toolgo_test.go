package toolgo

import (
	"github.com/lngwu11/toolgo/loggo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitLog(t *testing.T) {
	var logger = loggo.GetLogger("test")
	err := InitLog("debug", "")
	require.NoError(t, err)

	logger.Tracef("this is a trace log")
	logger.Debugf("this is a debug log")
	logger.Infof("this is a info log")
	logger.Warningf("this is a warning log")
	logger.Errorf("this is a error log")
}
