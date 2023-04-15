package log

import (
	"fmt"
	"testing"
)

func TestZap(t *testing.T) {

	type mdwware struct {
		Logger
	}

	logger, err := NewZapSugaredLogger(
		&Options{
			// MaxBackups: 1,
			MaxSize:  50 * 1024 * 1024,
			Compress: true,
			LogLevel: 3,
			LogPath:  "ms.log",
		})
	if err != nil {
		t.Fatal(err)
	}

	mw := mdwware{
		logger,
	}

	for i := 0; i < 10000; i++ {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			t.Parallel()
			mw.Logf(ERROR, "format: {%s}", "test error message")
		})
	}

	mw.Log(DEBUG, "test debug message")

}
