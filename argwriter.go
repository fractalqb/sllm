package sllm

import (
	"strconv"
	"time"
)

// ArgWriter is used by Expand to escape the argument when written as
// a value of a parameter. It is assumed that a user of this package
// should not use this type directly. However the type will be needed
// if one wants to provide an own ArgPrintFunc.
type ArgWriter []byte

// Write escapes the content so that it can be reliably recognized in a sllm
// message, i.e. replace each single backtick '`' with two backticks.
func (w *ArgWriter) Write(p []byte) (n int, err error) {
	n = len(*w)
	ep := 0
	for i, b := range p {
		if b == tmplEsc {
			*w = append(*w, p[ep:i+1]...)
			ep = i
		}
	}
	*w = append(*w, p[ep:]...)
	return len(*w) - n, nil
}

func (w *ArgWriter) WriteString(s string) (n int, err error) {
	n = len(*w)
	ep := 0
	for i, b := range []byte(s) {
		if b == tmplEsc {
			*w = append(*w, s[ep:i+1]...)
			ep = i
		}
	}
	*w = append(*w, s[ep:]...)
	return len(*w) - n, nil
}

func (w *ArgWriter) WriteInt(i int) (n int, err error) {
	l := len(*w)
	*w = itoa(*w, i, 0)
	return len(*w) - l, nil
}

func (w *ArgWriter) WriteInt64(i int64) (n int, err error) {
	l := len(*w)
	*w = strconv.AppendInt(*w, i, 10)
	return len(*w) - l, nil
}

func (w *ArgWriter) WriteTime(t time.Time, f TimeFormat) (n int, err error) {
	l := len(*w)
	*w = fmtTs(*w, t, f)
	return len(*w) - l, nil
}

type TimeFormat int

func (tf TimeFormat) Format(t time.Time) timeFormatter { return timeFormatter{tf, t} }

const (
	TUTC TimeFormat = 1 << iota
	Tdate
	Tweekday
	Tyear
	Tclock
	Tmillis
	Tmicros
)

const Tdefault = Tdate | Tweekday | Tclock

type timeFormatter struct {
	fmt TimeFormat
	t   time.Time
}

func fmtTs(buf []byte, t time.Time, fmt TimeFormat) []byte {
	if fmt&TUTC != 0 {
		t = t.UTC()
	}

	if fmt&(Tdate|Tyear|Tweekday) != 0 {
		ye, mo, dy := t.Date()
		if fmt&Tyear != 0 {
			buf = uitoa(buf, ye, 4)
			buf = append(buf, '-')
		}
		buf = uitoa(buf, int(mo), 2)
		buf = append(buf, '-')
		buf = uitoa(buf, dy, 2)
		if fmt&Tweekday != 0 {
			buf = append(buf, ' ')
			buf = append(buf, t.Weekday().String()[:2]...)
		}
		if fmt&(Tclock|Tmillis|Tmicros) != 0 {
			buf = append(buf, ' ')
		}
	}

	if fmt&(Tclock|Tmillis|Tmicros) != 0 {
		ho, mi, sc := t.Clock()
		buf = uitoa(buf, ho, 2)
		buf = append(buf, ':')
		buf = uitoa(buf, mi, 2)
		buf = append(buf, ':')
		buf = uitoa(buf, sc, 2)
		if fmt&Tmicros != 0 {
			buf = append(buf, '.')
			buf = uitoa(buf, t.Nanosecond()/1000, 6)
		} else if fmt&Tmillis != 0 {
			buf = append(buf, '.')
			buf = uitoa(buf, t.Nanosecond()/1000000, 3)
		}
	}
	return buf
}

func itoa(buf []byte, i, w int) []byte {
	if i < 0 {
		buf = append(buf, '-')
		i = -i
	}
	return uitoa(buf, i, w)
}

func uitoa(buf []byte, i, w int) []byte {
	var tmp = [20]byte{
		'0', '0', '0', '0', '0', '0', '0', '0', '0', '0',
		'0', '0', '0', '0', '0', '0', '0', '0', '0', '0',
	}
	wp := len(tmp) - 1
	w = len(tmp) - w
	for i > 9 {
		q := i / 10
		tmp[wp] = byte('0' + i - 10*q)
		wp--
		i = q
	}
	tmp[wp] = byte('0' + i)
	if wp > w {
		wp = w
	}
	return append(buf, tmp[wp:]...)
}
