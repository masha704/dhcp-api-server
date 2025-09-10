package logging

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

func (l *Logger) LogRequest(method, path, remoteAddr string, status int, duration time.Duration) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("%s - %s %s from %s - Status: %d - Duration: %v\n",
		timestamp, method, path, remoteAddr, status, duration)
	
	l.file.WriteString(message)
}

func (l *Logger) Close() {
	l.file.Close()
}