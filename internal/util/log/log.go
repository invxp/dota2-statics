package log

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

/*
工具包
日志库
*/

type Log struct {
	log *log.Logger
	mu  *sync.Mutex
}

func New(currentPath, pathName, fileName string, maxAge time.Duration, maxFileSize int64) *Log {
	if !filepath.IsAbs(pathName) {
		pathName = filepath.Join(currentPath, pathName)
	}

	fullFilePathName := filepath.Join(pathName, fileName)
	writer, err := rotatelogs.New(
		fullFilePathName+".%Y-%m-%d_%H:%M",
		rotatelogs.WithLinkName(fullFilePathName),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationSize(maxFileSize),
		rotatelogs.ForceNewFile(),
	)

	err = os.MkdirAll(pathName, 0755)
	if err != nil {
		log.Panic(err)
	}

	logFile, err := os.OpenFile(fullFilePathName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}

	return &Log{log.New(io.MultiWriter(logFile, writer, os.Stdout), "", log.LstdFlags|log.Lshortfile), &sync.Mutex{}}
}

func (l *Log) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.Printf(format, v...)
}
