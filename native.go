//go:build !js
// +build !js

package csl

import (
	"io/ioutil"
	"os"
)

var (
	Blue   = []byte("\033[0;36m")
	Green  = []byte("\033[1;32m")
	Red    = []byte("\033[1;31m")
	Purple = []byte("\033[0;35m")
	Brown  = []byte("\033[1;33m")
	Grey   = []byte("\033[0;37m")
	White  = []byte("\033[1;37m")
)

func StdLog() *Log {
	l := new(Log)

	l.Use(os.Stdout)

	l.err = os.Stderr

	return l
}

func (l *Log) Std() *Log {
	l.Use(os.Stdout)
	l.err = os.Stderr

	return l
}

func discardLog() *Log {
	l := new(Log)

	l.Use(ioutil.Discard)
	l.err = ioutil.Discard

	return l
}
