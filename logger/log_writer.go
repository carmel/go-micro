// 实现日志以文件形式保存，逻辑如下：
// 当前日志文件经LogPath指定，当文件大小达到MaxSize或文件时间跨度达到MaxAge，将当前日志文件压缩存档（存档名称经gzName方法获得），之后将LogPath文件清空

// 参考1: gopkg.in/natefinch/lumberjack.v2
// 参考2: https://github.com/easyCZ/logrotate
package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/carmel/go-micro/tool"
)

type LogWriter struct {
	ts   time.Time
	file *os.File
	dir  string
	opts Options
	size int64
	sync.Mutex
}

// NewLogWriter creates a new concurrency safe LogWriter which performs log rotation.
func NewLogWriter(opts Options) (*LogWriter, error) {
	dir := filepath.Dir(opts.LogPath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("directory %s does not exist and could not be created: %s", dir, err)
		}
	}

	// 以创建、只写、追加模式打开日志文件
	var f *os.File
	if f, err = os.OpenFile(opts.LogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644); err != nil {
		return nil, fmt.Errorf("failed to open file %s: %s", opts.LogPath, err)
	}

	w := &LogWriter{
		opts: opts,
		file: f,
		// bw:   bufio.NewWriter(f),
		dir: dir,
	}

	// w.bw.WriteString("\n")

	fi, _ := f.Stat()
	w.size = fi.Size() // + 1

	return w, nil
}

func (w *LogWriter) gzName() string {
	// return fmt.Sprintf("%s/%s-%s-%s.gz", w.dir, w.fn, time.Now().UTC().Format("2006-01-02T15-04-05.000"), util.RandomStr(3))
	return fmt.Sprintf("%s%c%s-%s.log.gz", w.dir, filepath.Separator, time.Now().UTC().Format(time.RFC3339), tool.RandomStr(3))
}

// Write implements io.Writer.
func (l *LogWriter) Write(p []byte) (n int, err error) {
	l.Lock()
	defer l.Unlock()

	writeLen := int64(len(p))
	if l.opts.MaxSize != 0 {
		if l.opts.MaxSize < writeLen {
			err = fmt.Errorf("write length %d exceeds maximum file size %d", writeLen, l.opts.MaxSize)
			return
		}

		if l.size+writeLen > l.opts.MaxSize {
			// 压缩存档
			err = l.archive()
			if err != nil {
				err = fmt.Errorf("failed to archive: %s", err)
				return
			}
		}
	}

	if l.opts.MaxAge != 0 && time.Now().After(l.ts.Add(l.opts.MaxAge)) {
		// 压缩存档
		err = l.archive()
		if err != nil {
			err = fmt.Errorf("failed to archive: %s", err)
			return
		}
	}

	n, err = l.file.Write(p)
	if err != nil {
		err = fmt.Errorf("failed to write to the log file: %s", err)
		return
	}
	l.size += writeLen
	return
}

// Close implements io.Closer
func (l *LogWriter) Close() error {
	l.Lock()
	defer l.Unlock()

	// err := l.bw.Flush()
	// if err != nil {
	// 	return fmt.Errorf("failed to flush buffered writer: %s", err)
	// }

	// File.Sync() 底层调用的是 fsync 系统调用，这会将数据和元数据都刷到磁盘
	// 如果只想刷数据到磁盘（比如，文件大小没变，只是变了文件数据），需要自己封装，调用 fdatasync（syscall.Fdatasync） 系统调用。
	// err = l.file.Sync()
	// if err != nil {
	// 	return fmt.Errorf("failed to sync current log file: %s", err)
	// }

	err := l.file.Close()
	if err != nil {
		return fmt.Errorf("failed to close current log file: %s", err)
	}

	l.file = nil

	return nil
}

func (l *LogWriter) archive() error {
	// 对 LogPath 指向的日志文件进行压缩存档 ===>
	gz := l.gzName()
	nf, err := os.Create(gz)
	if err != nil {
		return fmt.Errorf("failed to create gzip file: %s", err)
	}
	defer nf.Close()

	gzw := gzip.NewWriter(nf)
	if err != nil {
		return fmt.Errorf("failed to create gzip writer: %s", err)
	}
	defer gzw.Close()

	var f *os.File
	// 此处不能重用l.file，而要重新定义一个f，是因file是只写的
	f, err = os.Open(l.opts.LogPath)
	if err != nil {
		return fmt.Errorf("failed to open log file in read-only mode: %s", err)
	}
	defer f.Close()

	_, err = io.Copy(gzw, f)
	if err != nil {
		return fmt.Errorf("failed to write gzip file: %s", err)
	}

	// gzw.Flush()

	// 清空原来日志
	err = os.Truncate(l.opts.LogPath, 0)
	if err != nil {
		return fmt.Errorf("failed to truncate log file: %s", err)
	}

	// err = l.file.Truncate(0)
	// if err != nil {
	// 	return fmt.Errorf("failed to truncate log file: %s", err)
	// }
	// _, err = l.file.Seek(0, 0)
	// if err != nil {
	// 	return fmt.Errorf("failed to seek log file: %s", err)
	// }
	l.size = 0
	l.ts = time.Now().UTC()

	return nil
}
