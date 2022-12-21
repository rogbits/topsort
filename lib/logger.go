package lib

import (
	"log"
	"os"
)

type Logger struct {
	logger       *log.Logger
	TestMode     bool
	TestMessages []interface{}
}

func NewLogger() *Logger {
	l := new(Logger)
	l.logger = log.New(os.Stdout, "", log.Lshortfile)
	return l
}

func (l *Logger) Log(inputs ...interface{}) {
	if l.TestMode {
		l.TestMessages = append(l.TestMessages, inputs...)
		return
	}
	l.logger.Println(inputs...)
}

func (l *Logger) Fatal(inputs ...interface{}) {
	if l.TestMode {
		l.TestMessages = append(l.TestMessages, inputs...)
		return
	}
	l.logger.Fatal(inputs...)
}

func (l *Logger) Reset() {
	l.TestMessages = nil
}
