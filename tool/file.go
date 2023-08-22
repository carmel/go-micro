package tool

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const COMPRESSED_EXT = ".gz"

func FileGzip(f *os.File, gzName string) error {
	if gzName == "" {
		gzName = f.Name() + COMPRESSED_EXT
	}

	nf, err := os.Create(gzName)
	if err != nil {
		return err
	}
	defer nf.Close()

	w := gzip.NewWriter(nf)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, f)
	return err
	// 进行Sync读写时才需要使用Flush
	// return w.Flush()
}

func PathGzip(fpath, gzName string) error {
	var err error
	fpath, err = filepath.Abs(fpath)
	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	return FileGzip(f, gzName)
}

func PathCopy(srcPath, dstPath string) error {
	var err error
	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		return fmt.Errorf("srcPath error: %s", err)
	}

	dstPath, err = filepath.Abs(dstPath)
	if err != nil {
		return fmt.Errorf("dstPath error: %s", err)
	}

	var dat []byte
	dat, err = os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("read srcPath error: %s", err)
	}

	return os.WriteFile(dstPath, dat, 0644)
}

func FileCopy(srcFile *os.File, dstPath string) error {
	var err error
	dstPath, err = filepath.Abs(dstPath)
	if err != nil {
		return fmt.Errorf("dstPath error: %s", err)
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("creat dstPath: %s", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, srcFile)
	return err
}
