package csl

import "strconv"

type I8 int8

func (a *I8) Parse(f *Formatter) {
	v, _ := strconv.ParseInt(string(f.buf), 10, 8)
	*a = I8(v)
}

func (a I8) Fmt(f *Formatter) {
	f.buf = strconv.AppendInt(f.buf, int64(a), 10)
}

type I16 int16

func (a *I16) Parse(f *Formatter) {
	v, _ := strconv.ParseInt(string(f.buf), 10, 16)
	*a = I16(v)
}

func (a I16) Fmt(f *Formatter) {
	f.buf = strconv.AppendInt(f.buf, int64(a), 10)
}

type I32 int32

func (a *I32) Parse(f *Formatter) {
	v, _ := strconv.ParseInt(string(f.buf), 10, 32)
	*a = I32(v)
}

func (a I32) Fmt(f *Formatter) {
	f.buf = strconv.AppendInt(f.buf, int64(a), 10)
}

type I64 int64

func (a *I64) Parse(f *Formatter) {
	v, _ := strconv.ParseInt(string(f.buf), 10, 64)
	*a = I64(v)
}

func (a I64) Fmt(f *Formatter) {
	f.buf = strconv.AppendInt(f.buf, int64(a), 10)
}

//type Int int
//
//func (a *Int) Parse(f *Formatter) {
//	v, _ := strconv.ParseInt(string(f.buf), 10, 64)
//	*a = Int(v)
//}
//
//func (a Int) Fmt(f *Formatter) {
//	f.buf = strconv.AppendInt(f.buf, int64(a), 10)
//}

type U8 uint8

func (a *U8) Parse(f *Formatter) {
	v, _ := strconv.ParseUint(string(f.buf), 10, 8)
	*a = U8(v)
}

func (a U8) Fmt(f *Formatter) {
	f.buf = strconv.AppendUint(f.buf, uint64(a), 10)
}

type Byte byte

func (a *Byte) Parse(f *Formatter) {
	v, _ := strconv.ParseUint(string(f.buf), 10, 8)
	*a = Byte(v)
}

func (a Byte) Fmt(f *Formatter) {
	f.U8(byte(a))
}

type U16 uint16

func (a *U16) Parse(f *Formatter) {
	v, _ := strconv.ParseUint(string(f.buf), 10, 16)
	*a = U16(v)
}

func (a U16) Fmt(f *Formatter) {
	f.buf = strconv.AppendUint(f.buf, uint64(a), 10)
}

type U32 uint32

func (a *U32) Parse(f *Formatter) {
	v, _ := strconv.ParseUint(string(f.buf), 10, 32)
	*a = U32(v)
}

func (a U32) Fmt(f *Formatter) {
	f.buf = strconv.AppendUint(f.buf, uint64(a), 10)
}

type U64 uint64

func (a *U64) Parse(f *Formatter) {
	v, _ := strconv.ParseUint(string(f.buf), 10, 64)
	*a = U64(v)
}

func (a U64) Fmt(f *Formatter) {
	f.buf = strconv.AppendUint(f.buf, uint64(a), 10)
}

//type Uint uint
//
//func (a *Uint) Parse(f *Formatter) {
//	v, _ := strconv.ParseUint(string(f.buf), 10, 64)
//	*a = Uint(v)
//}
//
//func (a Uint) Fmt(f *Formatter) {
//	f.buf = strconv.AppendUint(f.buf, uint64(a), 10)
//}

type Rune rune

func (a *Rune) Parse(f *Formatter) {
	if len(f.buf) == 0 {
		return
	}

	*a = Rune([]rune(string(f.buf))[0])
}

func (a Rune) Fmt(f *Formatter) {
	f.Rune(rune(a))
}

type F32 float32

func (a *F32) Parse(f *Formatter) {
	v, _ := strconv.ParseFloat(string(f.buf), 32)
	*a = F32(v)
}

func (a F32) Fmt(f *Formatter) {
	f.validateFloatParams()
	f.buf = strconv.AppendFloat(f.buf, float64(a), f.FloatFmt, f.FloatPrec, 32)
}

type F64 float64

func (a *F64) Parse(f *Formatter) {
	v, _ := strconv.ParseFloat(string(f.buf), 64)
	*a = F64(v)
}

func (a F64) Fmt(f *Formatter) {
	f.validateFloatParams()
	f.buf = strconv.AppendFloat(f.buf, float64(a), f.FloatFmt, f.FloatPrec, 64)
}

//type F32f float32
//
//func (a F32f) Fmt(f *Formatter, fmt byte, prec int) {
//	f.FloatFmt = fmt
//	f.FloatPrec = prec
//	f.buf = strconv.AppendFloat(f.buf, float64(a), fmt, prec, 32)
//}
//
//type F64f float64
//
//func (a F64f) Fmt(f *Formatter, fmt byte, prec int) {
//	f.FloatFmt = fmt
//	f.FloatPrec = prec
//	f.buf = strconv.AppendFloat(f.buf, float64(a), fmt, prec, 64)
//}

type Bytes []byte

func (a *Bytes) Parse(f *Formatter) {
	copy(*a, f.buf)
}

func (a Bytes) Fmt(f *Formatter) {
	f.buf = append(f.buf, a...)
}

type Bool bool

func (a *Bool) Parse(f *Formatter) {
	v, _ := strconv.ParseBool(string(f.buf))
	*a = Bool(v)
}

func (a Bool) Fmt(f *Formatter) {
	f.buf = strconv.AppendBool(f.buf, bool(a))
}

type String string

func (a *String) Parse(f *Formatter) {
	*a = String(f.buf)
}

func (a String) Fmt(f *Formatter) {
	f.Str(string(a))
}

type I8s []int8
type I16s []int16
type I32s []int32
type I64s []int64
type Ints []int
type U8s []uint8
type U16s []uint16
type U32s []uint32
type U64s []uint64
type Uints []uint
type Runes []rune
type F32s []float32
type F64s []float64
type Bools []bool
type Strings []string

func (a I8s) Fmter(i int) Fmter {
	return (*I8)(&a[i])
}
func (a I16s) Fmter(i int) Fmter {
	return (*I16)(&a[i])
}
func (a I32s) Fmter(i int) Fmter {
	return (*I32)(&a[i])
}
func (a I64s) Fmter(i int) Fmter {
	return (*I64)(&a[i])
}

//func (a Ints) Fmter(i int) Fmter {
//	return (*Int)(&a[i])
//}
func (a U8s) Fmter(i int) Fmter {
	return (*U8)(&a[i])
}
func (a Bytes) Fmter(i int) Fmter {
	return (*Byte)(&a[i])
}
func (a U16s) Fmter(i int) Fmter {
	return (*U16)(&a[i])
}
func (a U32s) Fmter(i int) Fmter {
	return (*U32)(&a[i])
}
func (a U64s) Fmter(i int) Fmter {
	return (*U64)(&a[i])
}

//func (a Uints) Fmter(i int) Fmter {
//	return (*Uint)(&a[i])
//}
func (a Runes) Fmter(i int) Fmter {
	return (*Rune)(&a[i])
}
func (a F32s) Fmter(i int) Fmter {
	return (*F32)(&a[i])
}
func (a F64s) Fmter(i int) Fmter {
	return (*F64)(&a[i])
}

func (a Bools) Fmter(i int) Fmter {
	return (*Bool)(&a[i])
}

func (a Strings) Fmter(i int) Fmter {
	return (*String)(&a[i])
}

func (a *I8) Fmter(i int) Fmter {
	return a
}
func (a *I16) Fmter(i int) Fmter {
	return a
}
func (a *I32) Fmter(i int) Fmter {
	return a
}
func (a *I64) Fmter(i int) Fmter {
	return a
}

//func (a *Int) Fmter(i int) Fmter {
//	return a
//}
func (a *U8) Fmter(i int) Fmter {
	return a
}
func (a *Byte) Fmter(i int) Fmter {
	return a
}
func (a *U16) Fmter(i int) Fmter {
	return a
}
func (a *U32) Fmter(i int) Fmter {
	return a
}
func (a *U64) Fmter(i int) Fmter {
	return a
}

//func (a *Uint) Fmter(i int) Fmter {
//	return a
//}
func (a *Rune) Fmter(i int) Fmter {
	return a
}
func (a *F32) Fmter(i int) Fmter {
	return a
}
func (a *F64) Fmter(i int) Fmter {
	return a
}

func (a *Bool) Fmter(i int) Fmter {
	return a
}

func (a *String) Fmter(i int) Fmter {
	return a
}
