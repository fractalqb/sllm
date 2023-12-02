package sllm

import (
	"math"
	"time"
)

type TimeFormat int

func (tf TimeFormat) Fmt(t time.Time) timeFormatter { return timeFormatter{tf, t} }

const (
	TUTC TimeFormat = 1 << iota
	TNoDate
	TNoWeekday
	TYear
	TNoClock
	TMillis
	TMicros
)

const TDefault = TimeFormat(0)

type timeFormatter struct {
	fmt TimeFormat
	t   time.Time
}

func (tf timeFormatter) AppendSllm(buf []byte) []byte {
	if tf.fmt.anyOn(TUTC) {
		tf.t = tf.t.UTC()
	}
	needTime := tf.fmt.allOff(TNoClock)

	if tf.fmt.allOff(TNoDate) {
		ye, mo, dy := tf.t.Date()
		if tf.fmt.anyOn(TYear) {
			buf = uitoa(buf, ye, 4)
			buf = append(buf, '-')
		}
		buf = uitoa(buf, int(mo), 2)
		buf = append(buf, '-')
		buf = uitoa(buf, dy, 2)
		if tf.fmt.allOff(TNoWeekday) {
			buf = append(buf, ' ')
			buf = append(buf, tf.t.Weekday().String()[:2]...)
		}
		if needTime {
			buf = append(buf, ' ')
		}
	}

	if needTime {
		ho, mi, sc := tf.t.Clock()
		buf = uitoa(buf, ho, 2)
		buf = append(buf, ':')
		buf = uitoa(buf, mi, 2)
		buf = append(buf, ':')
		buf = uitoa(buf, sc, 2)
		switch {
		case tf.fmt.anyOn(TMicros):
			buf = append(buf, '.')
			buf = uitoa(buf, tf.t.Nanosecond()/1000, 6)
		case tf.fmt.anyOn(TMillis):
			buf = append(buf, '.')
			buf = uitoa(buf, tf.t.Nanosecond()/1000000, 3)
		}
		if tf.fmt.allOff(TUTC) {
			_, o := tf.t.Zone()
			buf = tzOff(buf, o, "+", "-")
		}
	} else if tf.fmt.allOff(TUTC) {
		_, o := tf.t.Zone()
		buf = tzOff(buf, o, " +", " -")
	}
	return buf
}

func tzOff(buf []byte, sec int, p, n string) []byte {
	sec = int(math.Round(float64(sec / (60 * 60))))
	switch {
	case sec < 0:
		buf = append(buf, n...)
		buf = uitoa(buf, -sec, 2)
	default:
		buf = append(buf, p...)
		buf = uitoa(buf, sec, 2)
	}
	return buf
}

// func itoa(buf []byte, i, w int) []byte {
// 	if i < 0 {
// 		buf = append(buf, '-')
// 		i = -i
// 	}
// 	return uitoa(buf, i, w)
// }

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

func (tf TimeFormat) anyOn(flags TimeFormat) (res bool) {
	return tf&flags != 0
}

func (tf TimeFormat) anyOff(flags TimeFormat) (res bool) {
	return tf&flags != flags
}

func (tf TimeFormat) allOn(flags TimeFormat) (res bool) {
	return tf&flags == flags
}

func (tf TimeFormat) allOff(flags TimeFormat) (res bool) {
	return tf&flags == 0
}
