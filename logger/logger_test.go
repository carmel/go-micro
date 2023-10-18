package logger

import (
	"sync"
	"testing"

	"github.com/carmel/go-micro/util"
)

type midware struct {
	Logger
}

func TestLog(t *testing.T) {
	Debugf("Debugf log: %s", "message")
	Infof("Infof log: %s", "message")
	Warnf("Warnf log: %s", "message")
	Errorf("Errorf log: %s", "message")
	Errorf("Errorf log message")
	With("att1", "22").Error("sffs")
	WithLog("att1", "22")(DEBUG, "sdfd")
}

func BenchmarkZap(b *testing.B) {

	logger, err := NewZapSugaredLogger(
		Options{
			// MaxBackups: 1,
			MaxSize: 10 * 1024 * 1024,
			// MaxSize:  10,
			LogLevel: DEBUG,
			LogPath:  "log/ms.log",
		},
	)
	if err != nil {
		b.Fatal(err)
	}

	mw := midware{
		logger,
	}

	wp := util.NewPool(10, &sync.WaitGroup{})

	b.StartTimer()

	for i := 0; i < 120; i++ {
		wp.Acquire()

		go func(i int) {
			defer wp.Release()
			mw.Logf(ERROR, "%d#format: {%s}", i, "test error message")
		}(i)

		// fmt.Printf("#%d\n", i+1)
	}
	wp.Wait()
	b.StopTimer()

	mw.Log(DEBUG, "test debug message")

}

func TestSlog(t *testing.T) {

	logger := NewSlogger(
		Options{
			// MaxBackups: 1,
			MaxSize: 10 * 1024 * 1024,
			// MaxSize:  10,
			LogLevel: DEBUG,
			LogPath:  "log/ms.log",
		})

	mw := midware{
		logger,
	}

	l := mw.With("author", "vector")
	l.Log(DEBUG, "debug msg")
	l.Log(INFO, "info msg")
	l.Log(WARN, "warn msg")
	l.Log(ERROR, "error msg")

	mw.Logf(DEBUG, "hello %s", "world")

}
