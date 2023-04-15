package log

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/microservices/util"
)

const (
	archiveTimeFormat = "2006-01-02T15-04-05.000"
)

var (
	// currentTime exists so it can be mocked out by tests.
	currentTime = time.Now
)

func NewLogWriter(opt *Options) (*LogWriter, error) {
	if opt.MaxSize <= 0 {
		return nil, errors.New("max size cannot be 0")
	}

	if opt.LogPath == "" {
		return nil, errors.New("log path cannot be empty")
	}

	r := &LogWriter{
		Options: *opt,
	}

	// 第一次打开日志文件，有则追加（重启服务时可能日志文件已存在），无则创建
	fl, err := os.OpenFile(r.LogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		// if we fail to open the old log file for some reason, just ignore
		// it and open a new log file.
		return nil, err
	}

	stat, _ := fl.Stat()
	r.file = fl
	r.size = stat.Size()

	return r, nil
}

type LogWriter struct {
	Options
	size      int64
	file      *os.File
	mu        sync.Mutex
	wg        sync.WaitGroup
	startMill sync.Once
}

func (r *LogWriter) Write(p []byte) (n int, err error) {
	writeLen := int64(len(p))
	if writeLen > r.MaxSize {
		return 0, fmt.Errorf(
			"write length %d, max size %d: write exceeds max file length", writeLen, r.MaxSize,
		)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 当要写入的内容大小大于限定值，则进行历史归档
	if r.MaxSize > 0 && r.size+writeLen >= r.MaxSize {
		r.wg.Add(1)
		r.startMill.Do(func() {
			go r.archive()
		})

		r.wg.Wait()
	} else {
		n, err = r.file.Write(p)
		r.size += int64(n)
	}

	return n, err
}

// Close implements io.Closer, and closes the current logfile.
func (r *LogWriter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.close()
}

// close closes the file if it is open.
func (r *LogWriter) close() error {
	if r.file == nil {
		return nil
	}
	err := r.file.Close()
	r.file = nil
	return err
}

// archive 进行历史归档
func (r *LogWriter) archive() {
	defer r.wg.Done()
	f, err := os.Open(r.LogPath)
	if err != nil {
		return
	}
	defer f.Close()

	// move the existing file
	filename := filepath.Base(r.LogPath)
	dir := filepath.Dir(filename)
	ext := filepath.Ext(filename)
	filename = fmt.Sprintf("%s-%s%s", filename[:len(filename)-len(ext)], currentTime().Format(archiveTimeFormat), ext)
	newname := filepath.Join(dir, filename)

	defer func() {
		// 如果报错则删除存档
		if err != nil {
			_ = os.Remove(newname)
		}
	}()
	if r.Compress {
		newname += util.COMPRESSED_EXT
		err = util.FileGzip(f, newname)
		if err != nil {
			return
		}
	} else {
		// 复制当前日志为存档文件
		err = util.FileCopy(f, newname)
		if err != nil {
			return
		}

	}

	if r.MaxAge > 0 {
		var files []fs.DirEntry
		files, err = os.ReadDir(filepath.Dir(r.LogPath))
		if err != nil {
			return
		}
		cutoff := currentTime().Add(-1 * r.MaxAge)
		var fn string
		// fi, _ := f.Stat()

		var fi os.FileInfo
		for _, f := range files {
			fn = f.Name()
			if f.IsDir() || fn == r.LogPath || fn == newname {
				continue
			}

			fi, _ = f.Info()
			if fi.ModTime().Before(cutoff) {
				_ = os.Remove(fi.Name())
			}
		}
	}

	// 清空日志
	err = os.Truncate(r.LogPath, 0)
	if err != nil {
		return
	}

	r.size = 0

}
