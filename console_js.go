//go:build js
// +build js

package csl

import console "honnef.co/go/js/console"

func Println(text string) {
	console.Log(text)
}

func Errorf(format string, err error) bool {
	if err != nil {
		console.Log(format, err.Error())
		return true
	}

	return false
}

func Printf(format string, a ...interface{}) {
	console.Log(format, a)
}

func Sprintf(format string, a ...interface{}) string {
	return "" // fmt.Sprintf(format+"\n", a...)
}
