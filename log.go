package csl

import (
	"io"
	"strconv"
	"time"
)

const AlwaysLogTime = false

type Log struct {
	out io.Writer
	err io.Writer

	errMode bool

	context []byte
	buffer  []byte
	depth   int

	Time bool
}

type discard struct{}

func (d discard) Write(b []byte) (int, error) {
	return len(b), nil
}

func (l *Log) Use(w io.Writer) {
	l.out = w
}

func (l *Log) Discard() *Log {
	l.out = discard{}
	l.err = discard{}

	return l
}

var (
	tab  = []byte("\t")
	tabs = []byte("\t\t\t\t\t\t\t\t\t\t\t")
	cl   = []byte(": ")
	nl   = []byte("\n")
)

func (l *Log) SetContext(context string) *Log {
	l.context = []byte(context)

	return l
}

func (l *Log) AppendContext(context string) {
	l.context = append(l.context, cl...)
	l.context = append(l.context, context...)
}

func (l Log) With() *Log {
	l.buffer = nil
	l.depth = 0

	return &l
}

func (l *Log) Indent() *Log {
	l.depth++
	if l.depth > 10 {
		l.depth = 10
	}
	return l
}

func (l *Log) Outdent() *Log {
	l.depth--
	if l.depth < 0 {
		l.depth = 0
	}
	return l
}

func (l *Log) tabs(offset int) {
	l.buffer = append(l.buffer, tabs[0:(l.depth+offset)]...)
}

func (l *Log) cl() { l.buffer = append(l.buffer, cl...) }
func (l *Log) nl() { l.buffer = append(l.buffer, nl...) }

func (l *Log) F() {
	l.Flush()
}

func (l *Log) FatalIf(msg string, context string, err error) bool {
	if err != nil {
		l.Error(msg).
			Err(context, err).Fatal()
		return true
	}

	return false
}

func (l *Log) ErrorIf(msg string, context string, err error) bool {
	return l.IfError(msg, context, err)
}

func (l *Log) IfError(msg string, context string, err error) bool {
	if err != nil {
		l.Error(msg).
			Err(context, err).F()
		return true
	}

	return false
}

func (l *Log) Fatal() {
	l.errMode = true
	pErr := string(l.buffer)
	l.Flush()

	panic(pErr)
}

func (l *Log) Fatalf(fmt string, a ...interface{}) {
	l.Error("fatal error").Printf("err", fmt, a...).Fatal()
}

func (l *Log) Flush() {
	if l.out == nil || l.err == nil {
		return
	}

	if l.Time || AlwaysLogTime {
		l.AppendTime()
	}

	l.Clr(White)

	if !l.errMode {
		l.out.Write(l.buffer)
	} else {
		l.err.Write(l.buffer)
	}

	l.buffer = l.buffer[:0]

	l.errMode = false

	//l.buffers.Push(l.buffer)
}

func (l *Log) AppendTime() *Log {
	return l.String("time", time.Now().Format("15:04:05.000000000"))
}

func (l *Log) T() {
	if !l.Time && !AlwaysLogTime {
		l.AppendTime()
	}

	l.Flush()
}

func (l *Log) Clr(ansi []byte) *Log {
	l.buffer = append(l.buffer, ansi...)
	return l
}

func (l *Log) Err(msg string, err error) *Log {
	if err == nil {
		return l
	}

	l.errMode = true

	return l.Clr(Red).String(msg, err.Error()).Clr(White)
}

func (l *Log) Error(msg string) *Log {
	l.errMode = true
	//l.buffer, _ = l.buffers.Pop()

	return l.Clr(Red).info(msg).Clr(White)
}

func (l *Log) Success(msg string) *Log {
	return l.Clr(Green).info(msg).Clr(White)
}

func (l *Log) Notice(msg string) *Log {
	return l.Clr(Purple).info(msg).Clr(White)
}

func (l *Log) Warn(msg string) *Log {
	return l.Clr(Brown).info(msg).Clr(White)
}

func (l *Log) giveContext() {
	if l.context != nil {
		l.bytes(l.context)
		l.cl()
	}
}

func (l *Log) info(msg string) *Log {
	l.tabs(0)
	l.Clr(White)
	l.giveContext()
	l.string(msg)
	//l.cl()
	l.nl()
	return l
}

func (l *Log) Info(msg string) *Log {
	l.errMode = false
	//l.buffer, _ = l.buffers.Pop()

	return l.info(msg)
}

func (l *Log) desc(desc string) {
	l.tabs(1)
	l.string(desc)
	l.cl()
}

func (l *Log) Val(name string, val interface{}) *Log {
	return l.String(name, Sprintf("%v", val))
}

func (l *Log) Struct(name string, val interface{}) *Log {
	return l.String(name, Sprintf("%+v", val))
}

func (l *Log) Printf(name, fmt string, a ...interface{}) *Log {
	return l.String(name, Sprintf(fmt, a...))
}

func (l *Log) Errorf(fmt string, a ...interface{}) *Log {
	return l.String("", Sprintf(fmt, a...))
}

func (l *Log) bytes(bytes []byte) {
	l.buffer = append(l.buffer, bytes...)
}

func (l *Log) Bytes(name string, bytes []byte) *Log {
	l.tabs(1)
	l.string(name)
	l.cl()
	l.bytes(bytes)
	l.nl()

	return l
}

func (l *Log) string(str string) {
	l.buffer = append(l.buffer, str...)
}

func (l *Log) String(name, str string) *Log {
	l.tabs(1)
	l.string(name)
	l.cl()
	l.string(str)
	l.nl()
	return l
}

func (l *Log) quote(str string) {
	l.buffer = strconv.AppendQuote(l.buffer, str)
}

func (l *Log) Quote(name, str string) *Log {
	l.tabs(1)
	l.string(name)
	l.cl()
	l.quote(str)
	l.nl()
	return l
}

type Stringer interface {
	String() string
}

func (l *Log) Stringer(name string, str Stringer) *Log {
	l.String(name, str.String())
	return l
}

func (l *Log) QuoteStringer(name string, str Stringer) *Log {
	l.Quote(name, str.String())
	return l
}

func (l *Log) i64(i int64) {
	l.buffer = strconv.AppendInt(l.buffer, i, 10)
}

func (l *Log) I64(desc string, i int64) *Log {
	l.desc(desc)
	l.i64(i)
	l.nl()
	return l
}

func (l *Log) Int(desc string, i int) *Log {
	l.I64(desc, int64(i))
	return l
}

func (l *Log) I32(desc string, i int32) *Log {
	l.I64(desc, int64(i))
	return l
}

func (l *Log) I16(desc string, i int16) *Log {
	l.I64(desc, int64(i))
	return l
}

func (l *Log) I8(desc string, i int8) *Log {
	l.I64(desc, int64(i))
	return l
}

func (l *Log) u64(u uint64) {
	l.buffer = strconv.AppendUint(l.buffer, u, 10)
}

func (l *Log) Uint(desc string, u uint) *Log {
	l.U64(desc, uint64(u))
	return l
}

func (l *Log) U64(desc string, u uint64) *Log {
	l.desc(desc)
	l.u64(u)
	l.nl()
	return l
}

func (l *Log) U32(desc string, u uint32) *Log {
	l.U64(desc, uint64(u))
	return l
}

func (l *Log) U16(desc string, u uint16) *Log {
	l.U64(desc, uint64(u))
	return l
}

func (l *Log) U8(desc string, u uint8) *Log {
	l.U64(desc, uint64(u))
	return l
}

func (l *Log) Byte(desc string, b byte) *Log {
	l.U8(desc, b)
	return l
}

func (l *Log) bool(b bool) {
	l.buffer = strconv.AppendBool(l.buffer, b)
}

func (l *Log) Bool(desc string, b bool) *Log {
	l.desc(desc)
	l.bool(b)
	l.nl()
	return l
}

func (l *Log) f64(f float64) {
	l.buffer = strconv.AppendFloat(l.buffer, f, 'g', -1, 64)
}

func (l *Log) F64(desc string, f float64) *Log {
	l.desc(desc)
	l.f64(f)
	l.nl()
	return l
}

func (l *Log) f32(f float32) {
	l.buffer = strconv.AppendFloat(l.buffer, float64(f), 'g', -1, 32)
}

func (l *Log) F32(desc string, f float32) *Log {
	l.desc(desc)
	l.f32(f)
	l.nl()

	return l
}

func (l *Log) Auto(msg string, fields ...interface{}) {
	l.tabs(0)
	l.string(msg)
	l.nl()

	for i := 0; i < len(fields); i++ {
		if i%2 == 0 {
			l.tabs(1)
		}

		switch t := fields[i].(type) {
		case string:
			l.string(t)
		case int:
			l.i64(int64(t))
		}

		if i%2 == 0 {
			l.cl()
		} else {
			l.nl()
		}
	}

	l.Flush()
}
