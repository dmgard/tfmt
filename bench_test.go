package tfmt

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func BenchmarkFmt1Int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fmt.Fprintln(ioutil.Discard, "test", "test", 23)
	}
}

func BenchmarkFmt10Int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fmt.Fprintln(ioutil.Discard, "test", "test", 23, "test", 23, "test", 23, "test", 23, "test", 23,
			"test", 23, "test", 23, "test", 23, "test", 23,
			"test", 23, "test", 23)
	}
}
func BenchmarkFmt3Strings(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fmt.Fprintln(ioutil.Discard, "test", "test", "test")
	}
}
