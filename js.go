// +build js

package csl

import (
	"honnef.co/go/js/console"
)

var (
	Blue   = []byte{}
	Green  = []byte{}
	Red    = []byte{}
	Purple = []byte{}
	Brown  = []byte{}
	Grey   = []byte{}
	White  = []byte{}
)

var std = new(stdWriter)

type stdWriter struct{}

func (s *stdWriter) Write(buf []byte) (int, error) {
	console.Log(string(buf))

	return len(buf), nil
}

func StdLog() *Log {
	l := new(Log)

	l.Use(std)
	l.err = std

	return l
}

func (l *Log) Std() *Log {
	l.Use(std)
	l.err = std

	return l
}

func discardLog() *Log {
	l := new(Log)

	l.Use(discard{})
	l.err = discard{}

	return l
}
