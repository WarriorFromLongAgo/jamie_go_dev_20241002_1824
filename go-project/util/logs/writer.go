package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var myLogWriter dateFileWriter

var targetFile *os.File

var logTime string

type dateFileWriter struct {
	io.Writer
}

func init() {
	RefreshLogFile()
	log.SetOutput(io.MultiWriter(os.Stdout, &myLogWriter))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func (d *dateFileWriter) Write(p []byte) (n int, err error) {
	if logTime != time.Now().Format(LogTimeFormat) {
		RefreshLogFile()
	}
	return targetFile.Write(p)
}

func RefreshLogFile() {
	logTime = time.Now().Format(LogTimeFormat)
	f, err := GetLogFile()
	if err != nil {
		fmt.Println(err)
		return
	}
	targetFile = f
	fmt.Println("RefreshLogFile success~")
}

func MyLogWriter() io.Writer {
	return io.MultiWriter(os.Stdout, &myLogWriter)
}
