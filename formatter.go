package csl

import (
	"bytes"
	"strconv"
	
	"golang.org/x/exp/constraints"
)

var tabsS = "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"

func Fmt() *Formatter {
	return &Formatter{}
}

type Fmter interface {
	Fmt(*Formatter)
	Parse(*Formatter)
}

type Fmters interface {
	Fmter(int) Fmter
}

type FmterSlice []Fmter

func (f FmterSlice) Fmter(i int) Fmter {
	return f[i]
}

type Formatter struct {
	buf       []byte
	FloatFmt  byte
	FloatPrec int
	tab       int
	space int
	ifs       []int // record of positions to rewind to if a call to "fi" is encountered with false

	pos int // last printed location
}

type ErrFmt Formatter

func (f ErrFmt) Error() string {
	return Formatter(f).String()
}

func (f *Formatter) Write(b []byte) (n int, err error) {
	f.Bytes(b)
	return len(b), nil
}

type FmtMaybe []*Formatter

func (f FmtMaybe) MustUnwrap() *Formatter {
	if len(f) == 0 {
		return &Formatter{}
	}

	if len(f) == 1 {
		if f[0] == nil {
			panic("do not pass a nil *Formatter to MustUnwrap")
		}
		return f[0]
	}

	panic("don't put more than one formatter into a FmtMaybe")
}

func (f *Formatter) Clear() *Formatter {
	f.buf = f.buf[:0]
	f.tab = 0
	return f
}

func (f *Formatter) Free() {
	f.buf = nil
}

func (f Formatter) String() string {
	return string(f.buf)
}

func (f *Formatter) YetUnseen() string {
	pos := f.pos
	f.pos = len(f.buf)
	return string(f.buf[pos:])
}

func (f *Formatter) LastN(n int) string {
	if n > len(f.buf) {
		n = len(f.buf)
	}

	if n < 0 {
		n = 0
	}
	return string(f.buf[len(f.buf)-n:])
}

func (f *Formatter) Buf() []byte {
	return f.buf
}

func (f *Formatter) Equals(g *Formatter) bool {
	if len(f.buf) != len(g.buf) {
		return false
	}

	for i := range f.buf {
		if f.buf[i] != g.buf[i] {
			return false
		}
	}

	return true
}

func (f *Formatter) Runes() []rune {
	return []rune(string(f.buf))
}

func (f *Formatter) SetBuf(buf []byte) {
	f.buf = buf
}

func (f *Formatter) Clone(g *Formatter) {
	if cap(f.buf) < len(g.buf) {
		f.buf = make([]byte, len(g.buf))
	}

	f.buf = f.buf[:len(g.buf)]

	copy(f.buf, g.buf)
}

func (f Formatter) Pos() int {
	return len(f.buf)
}

func (f Formatter) Rewound(to int) *Formatter {
	f.buf = f.buf[:to]
	return &f
}

func (f *Formatter) Rewind(to int) *Formatter {
	f.buf = f.buf[:to]
	return f
}

// TODO Go's builtin strcmp has gotta be orders of magnitude faster than this - if string([]byte) is done without allocations it's a no brainer
func (f *Formatter) ReplaceIfDiff(fm Fmter) bool {
	src := f.buf

	fm.Fmt(f)
	nw := f.buf[len(src):]

	diff := false

	for i := range nw {
		diff = diff || f.buf[i] == nw[i]
		f.buf[i] = nw[i]
	}

	f.buf = f.buf[:len(nw)]
	return diff
}

func (f *Formatter) Quote(a string) *Formatter {
	f.buf = strconv.AppendQuote(f.buf, a)
	return f
}

func (f *Formatter) RawQuote(a string) *Formatter {
	return f.Str("`").Str(a).Str("`")
}

func (f *Formatter) Str(a string) *Formatter {
	f.buf = append(f.buf, a...)
	return f
}

func (f *Formatter) Printf(fmt string, a ...interface{}) *Formatter {
	return f.Str(Sprintf(fmt, a...))
}

func (f *Formatter) Len() int {
	return len(f.buf)
}

func (f *Formatter) Delete(at, len int) *Formatter {
	f.buf = append(f.buf[:at], f.buf[at+len:]...)
	return f
}

func (f *Formatter) InsertBytes(at int, b []byte) *Formatter {
	if at >= len(f.buf) {
		f.buf = append(f.buf, b...)
		return f
	}

	f.buf = append(f.buf[:at+len(b)], f.buf[at:]...)
	copy(f.buf[at:], b)
	return f
}

func (f *Formatter) InsertString(at int, b string) *Formatter {
	if at >= len(f.buf) {
		f.buf = append(f.buf, b...)
		return f
	}

	f.buf = append(f.buf[:at+len(b)], f.buf[at:]...)
	copy(f.buf[at:], b)
	return f
}

func (f *Formatter) ReplaceBytes(start, end int, b []byte) *Formatter {
	f.buf = append(f.buf[:start+len(b)], f.buf[end:]...)
	copy(f.buf[start:], b)
	return f
}

func (f *Formatter) ReplaceString(start, end int, b string) *Formatter {
	if start+len(b) > len(f.buf) {
		buf := make([]byte, len(f.buf)-(end-start)+len(b))
		copy(buf, f.buf[:start])
		copy(buf[start:], b)
		copy(buf[start+len(b):], f.buf[end:])
		f.buf = buf
		return f
	}

	f.buf = append(f.buf[:start+len(b)], f.buf[end:]...)
	copy(f.buf[start:], b)
	return f
}

func (f *Formatter) ContainsRune(r rune) bool {
	return bytes.ContainsRune(f.buf, r)
}

func (f *Formatter) Tab() *Formatter {
	f.buf = append(f.buf, []byte("\t")...)
	return f
}

func (f *Formatter) Indent() *Formatter {
	return f.IndentTabs(1)
}

func (f *Formatter) Outdent() *Formatter {
	return f.OutdentTabs(1)
}

func (f *Formatter) IndentTabs(n int) *Formatter {
	f.tab = max(f.tab+n, 0)
	return f.Tabs(n)
}

func (f *Formatter) OutdentTabs(n int) *Formatter {
	f.tab = max(f.tab-n, 0)
	return f.TrimRightCopiesN("\t", n)
}

func (f *Formatter) TrimRightCopiesN(s string, n int) *Formatter {
	for pos := len(f.buf) - len(s);
		n > 0 && pos >= 0 && string(f.buf[pos:]) == s;
	pos, n = len(f.buf)-len(s), n-1 {
		f.pos -= len(s)
		f.buf = f.buf[:pos]
	}
	return f
}

func (f *Formatter) IndentSpaces(n int) *Formatter {
	f.space = max(f.space+n, 0)
	return f.Spaces(n)
}

func (f *Formatter) OutdentSpaces(n int) *Formatter {
	f.space = max(f.space-n, 0)
	return f.TrimRightCopiesN(" ", n)
}

const spaces = "                                                "

func (f *Formatter) Spaces(n int) *Formatter {
	for n > len(spaces) {
		f.Str(spaces)
		n -= len(spaces)
	}
	f.Str(spaces[:n])
	return f
}

func Tabs[I constraints.Integer](n I) string {
	return (&Formatter{}).Tabs(int(n)).String()
}

func (f *Formatter) Tabs(n int) *Formatter {
	for n > len(tabsS) {
		f.Str(tabsS)
		n -= len(tabsS)
	}
	f.Str(tabsS[:n])
	return f
}

func (f *Formatter) putIndent() *Formatter {
	return f.Tabs(f.tab).Spaces(f.space)
}

// Ln TODO optionally take a number of lines
func (f *Formatter) Ln() *Formatter {
	f.buf = append(f.buf, []byte("\n")...)
	return f.putIndent()
}

func (f *Formatter) Cln() *Formatter {
	f.buf = append(f.buf, []byte(",\n")...)
	return f.putIndent()
}

func Int[T ~int | ~int8 | ~int16 | ~int32 | ~int64](f *Formatter, i T) *Formatter {
	return f.I64(int64(i))
}
func Uint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](f *Formatter, i T) *Formatter {
	return f.U64(uint64(i))
}

func (f *Formatter) I8(a int8) *Formatter {
	f.buf = strconv.AppendInt(f.buf, int64(a), 10)
	return f
}

func (f *Formatter) I16(a int16) *Formatter {
	f.buf = strconv.AppendInt(f.buf, int64(a), 10)
	return f
}

func (f *Formatter) I32(a int32) *Formatter {
	f.buf = strconv.AppendInt(f.buf, int64(a), 10)
	return f
}

func (f *Formatter) I64(a int64) *Formatter {
	f.buf = strconv.AppendInt(f.buf, int64(a), 10)
	return f
}

func (f *Formatter) Int(a int) *Formatter {
	f.buf = strconv.AppendInt(f.buf, int64(a), 10)
	return f
}

func (f *Formatter) U8(a uint8) *Formatter {
	f.buf = strconv.AppendUint(f.buf, uint64(a), 10)
	return f
}

func (f *Formatter) Byte(a byte) *Formatter {
	return f.U8(a)
}

func (f *Formatter) U16(a uint16) *Formatter {
	f.buf = strconv.AppendUint(f.buf, uint64(a), 10)
	return f
}

func (f *Formatter) Stringer(a Stringer) *Formatter {
	return f.Str(a.String())
}

func (f *Formatter) U32(a uint32) *Formatter {
	f.buf = strconv.AppendUint(f.buf, uint64(a), 10)
	return f
}

func (f *Formatter) U64(a uint64) *Formatter {
	f.buf = strconv.AppendUint(f.buf, uint64(a), 10)
	return f
}

func (f *Formatter) Uint(a uint) *Formatter {
	f.buf = strconv.AppendUint(f.buf, uint64(a), 10)
	return f
}

func (f *Formatter) Rune(a rune) *Formatter {
	return f.Str(string(a))
}

func (f *Formatter) RuneQuote(a rune) *Formatter {
	f.buf = strconv.AppendQuoteRune(f.buf, a)
	return f
}

func (f *Formatter) validateFloatParams() {
	if f.FloatFmt == 0 {
		f.FloatFmt = 'g'
	}
	if f.FloatPrec == 0 {
		f.FloatPrec = -1
	}
}

func (f *Formatter) F32(a float32) *Formatter {
	f.validateFloatParams()
	f.buf = strconv.AppendFloat(f.buf, float64(a), f.FloatFmt, f.FloatPrec, 32)
	return f
}

func (f *Formatter) F64(a float64) *Formatter {
	f.validateFloatParams()
	f.buf = strconv.AppendFloat(f.buf, float64(a), f.FloatFmt, f.FloatPrec, 64)
	return f
}

func (f *Formatter) F32f(a float32, fmt byte, prec int) *Formatter {
	f.FloatFmt = fmt
	f.FloatPrec = prec
	f.buf = strconv.AppendFloat(f.buf, float64(a), fmt, prec, 32)
	return f
}

func (f *Formatter) F64f(a float64, fmt byte, prec int) *Formatter {
	f.FloatFmt = fmt
	f.FloatPrec = prec
	f.buf = strconv.AppendFloat(f.buf, float64(a), fmt, prec, 64)
	return f
}

func (f *Formatter) Bytes(a []byte) *Formatter {
	f.buf = append(f.buf, a...)
	return f
}

func (f *Formatter) Bool(a bool) *Formatter {
	f.buf = strconv.AppendBool(f.buf, a)
	return f
}

func (f *Formatter) TrimRight(s string) *Formatter {
	f.buf = bytes.TrimRight(f.buf, s)

	return f
}

func (f *Formatter) SetIndent(depth int) *Formatter {
	f.tab = depth

	if f.tab > len(tabsS) {
		f.tab = len(tabsS)
	} else if f.tab < 0 {
		f.tab = 0
	}

	return f
}

func (f *Formatter) IndentTo(depth int) *Formatter {
	return f.SetIndent(depth).TrimRight("\t").putIndent()
}

func (f *Formatter) LnIfIndent(depth int) *Formatter {
	tb := f.tab
	f.SetIndent(depth)
	if tb != depth {
		f.Ln()
	}

	return f
}

func (f *Formatter) Fatalf(fmt string, a ...interface{}) {
	f.Printf(fmt, a...)
	panic(f.String())
}

func (f *Formatter) Println(s string) *Formatter {
	return f.Str(s).Ln()
}

func (f *Formatter) OutdentLn() *Formatter {
	return f.TrimRight("\n\t").Outdent().Ln()
}

func (f *Formatter) RuneLit(r rune) *Formatter {
	return f.Str(strconv.QuoteRune(r))
}

type FmtRecvr interface {
	Fmt(*Formatter) *Formatter
}

type FmtMaybeRecvr interface {
	Fmt(...*Formatter) *Formatter
}

func (f *Formatter) Fmter(r interface{}) *Formatter {
	switch t := r.(type) {
	case FmtRecvr:
		return t.Fmt(f)
	case FmtMaybeRecvr:
		return t.Fmt(f)
	default:
		panic("invalid type for Fmter")
	}
}

func (f *Formatter) If() *Formatter {
	f.ifs = append(f.ifs, f.Pos())
	return f
}
func (f *Formatter) Fi(b bool) *Formatter {
	if len(f.ifs) == 0 {
		panic("Fi called in formatter without matching If")
	}
	if !b {
		f.Rewind(f.ifs[len(f.ifs)-1])
	}
	f.ifs = f.ifs[:len(f.ifs)-1]
	return f
}

func (f *Formatter) Intf(i int, base int) *Formatter {
	return f.Str(strconv.FormatInt(int64(i), base))
}

func (f *Formatter) Reset() {
	f.Rewind(0)
}
