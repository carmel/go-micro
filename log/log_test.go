package log

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/microservices/util"
)

type mdware struct {
	Logger
}

func BenchmarkZap(b *testing.B) {

	logger, err := NewZapSugaredLogger(
		Options{
			// MaxBackups: 1,
			MaxSize: 10 * 1024 * 1024,
			// MaxSize:  10,
			LogLevel: 0,
			LogPath:  "log/ms.log",
		},
	)
	if err != nil {
		b.Fatal(err)
	}

	mw := mdware{
		logger,
	}

	wp := util.NewPool(10, &sync.WaitGroup{})

	b.StartTimer()

	for i := 0; i < 120000; i++ {
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

func TestOnce(t *testing.T) {

	doSomething := func() {
		fmt.Println("do something only once...")
	}

	once := &sync.Once{}
	wg := &sync.WaitGroup{}

	var routines int = 4
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Printf("goroutine %d...\n", i+1)
			once.Do(doSomething)
		}(i)
	}

	wg.Wait()
}

func TestAny(t *testing.T) {
	fmt.Println(len("\n"))

	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	bw := bufio.NewWriter(f)

	bw.WriteString("\n")
	// bw.Flush()

	fi, _ := f.Stat()
	fmt.Println(fi.Size())
}
