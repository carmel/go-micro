package tool

import (
	"fmt"
	"os"
	"path/filepath"
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

func TestPool2(t *testing.T) {

}
