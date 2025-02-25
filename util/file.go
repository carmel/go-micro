package util

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
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

func FileZip(writer io.Writer, compressMethod uint16, src ...struct {
	Name string
	Path string
}) error {
	zw := zip.NewWriter(writer)
	defer zw.Close()

	var (
		err    error
		fl     *os.File
		nw     io.Writer
		info   fs.FileInfo
		header *zip.FileHeader
	)
	for _, v := range src {
		fl, err = os.Open(v.Path)
		if err != nil {
			return err
		}
		defer fl.Close()
		info, err = fl.Stat()
		if err != nil {
			return err
		}
		header, err = zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		// 压缩算法
		header.Method = compressMethod
		// header.Name = filepath.Base(v)
		header.Name = v.Name
		nw, err = zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(nw, fl)
		if err != nil {
			return err
		}
	}
	return nil
}
