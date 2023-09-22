package tool

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestCompress(t *testing.T) {
	p := "../README.md"
	p, _ = filepath.Abs(p)
	// f, _ := os.Open(p)
	// defer f.Close()
	// fmt.Println(FileGzip(f))

	// fmt.Println(PathGzip(p))

	dat, err := os.ReadFile(p)
	fmt.Println(string(dat), err)

	fmt.Println(os.FileMode(0644).String())
}

func TestPool(t *testing.T) {

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
