package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"time"
)

var (
	logName string

	logInterval time.Duration
)

func main() {

	flag.StringVar(&logName, "log", "logs/main.log", "log path and filename")
	flag.DurationVar(&logInterval, "interval", time.Second, "log interval")

	flag.Parse()

	logger := NewLogger(logName)

	for {
		err := logger.Info(fmt.Sprintf("Random number: %d", rand.Intn(100)))
		if err != nil {
			log.Printf("Failed to write log: %v", err)
		}

		time.Sleep(logInterval)
	}
}

type Logger struct {
	filename string
	file     *os.File
	count    int
}

func NewLogger(filename string) *Logger {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		pathName := path.Dir(filename)
		os.Mkdir(pathName, 0777)
		_, err = os.Create(filename)
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	return &Logger{
		filename: filename,
		file:     file,
	}
}

func (l *Logger) Info(msg string) error {
	logEntry := struct {
		Timestamp time.Time `json:"timestamp"`
		Message   string    `json:"message"`
	}{
		Timestamp: time.Now(),
		Message:   msg,
	}

	jsonEntry, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %v", err)
	}

	if l.count%5 == 0 {
		fmt.Printf("[%s] %s\n", logEntry.Timestamp.Format("2006-01-02 15:04:05.006"), logEntry.Message)
	}

	_, err = l.file.Write(jsonEntry)
	l.file.WriteString("\n")
	if err != nil {
		return fmt.Errorf("failed to write log entry to file: %v", err)
	}

	l.count++

	if l.count%4 == 0 {
		fileInfo, err := os.Stat(l.filename)
		if err != nil {
			return fmt.Errorf("failed to get file info: %v", err)
		}

		fileSize := fileInfo.Size()
		if fileSize >= 1024 {
			err = l.rotate()
			if err != nil {
				return fmt.Errorf("failed to rotate log file: %v", err)
			}
		}
	}

	return nil
}

func (l *Logger) rotate() error {
	err := l.file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %v", err)
	}

	now := time.Now()
	newFilename := fmt.Sprintf("%s.%04d-%02d-%02d_%02d-%02d", l.filename, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	err = os.Rename(l.filename, newFilename)
	if err != nil {
		return fmt.Errorf("failed to rename file: %v", err)
	}

	l.file, err = os.Create(l.filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	return nil
}
