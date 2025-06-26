//go:build !js
// +build !js

package csl

import "fmt"

func Println(text string) {
	fmt.Println(text)
}

func Errorf(format string, err error) bool {
	if err != nil {
		fmt.Printf(format, err)
		return true
	}

	return false
}

func Printf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func Printf2D(format string, width, height int, a ...interface{}) {
	rows := len(a) / width

	var i, start, end int

	for i = 0; i < rows; i++ {
		end += width

		if end > len(a) {
			end = len(a)
		}

		Printf(format, a[start:end])

		if end == len(a) {
			break
		}
	}

	for i < height {
		Printf("")
		i++
	}
}
