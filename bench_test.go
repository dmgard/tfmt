package csl

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/uber-go/zap"
)

func BenchmarkStdLog(b *testing.B) {
	for n := 0; n < b.N; n++ {
		StdLog()
	}
}

func BenchmarkWrite1Int(b *testing.B) {
	l := discardLog()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		l.Info("test").Int("test", 33).Flush()
	}
}

func BenchmarkZap1Int(b *testing.B) {
	l := zap.New(zap.NewJSONEncoder())
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		l.Debug("test", zap.Int("test", 23))
	}
}

func BenchmarkFmt1Int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fmt.Fprintln(ioutil.Discard, "test", "test", 23)
	}
}

func BenchmarkWrite10Int(b *testing.B) {
	l := discardLog()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		l.Info("test").
			Int("test", 33).Int("test", 33).Int("test", 33).
			Int("test", 33).Int("test", 33).Int("test", 33).
			Int("test", 33).Int("test", 33).Int("test", 33).
			Int("test", 33).
			Flush()
	}
}

func BenchmarkZap10Int(b *testing.B) {
	l := zap.New(zap.NewJSONEncoder())
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		l.Debug("test", zap.Int("test", 23), zap.Int("test", 23), zap.Int("test", 23), zap.Int("test", 23), zap.Int("test", 23), zap.Int("test", 23), zap.Int("test", 23), zap.Int("test", 23), zap.Int("test", 23), zap.Int("test", 23))
	}
}

func BenchmarkFmt10Int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fmt.Fprintln(ioutil.Discard, "test", "test", 23, "test", 23, "test", 23, "test", 23, "test", 23,
			"test", 23, "test", 23, "test", 23, "test", 23,
			"test", 23, "test", 23)
	}
}

func BenchmarkZap3Strings(b *testing.B) {
	l := zap.New(zap.NewJSONEncoder())
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		l.Debug("test", zap.String("test", "test"))
	}
}

func BenchmarkWrite3Strings(b *testing.B) {
	l := discardLog()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		l.Info("test").String("test", "test").Flush()
	}
}

func BenchmarkFmt3Strings(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fmt.Fprintln(ioutil.Discard, "test", "test", "test")
	}
}
