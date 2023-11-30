package sllm

import "time"

type TimeFormat int

func (tf TimeFormat) Fmt(t time.Time) timeFormatter { return timeFormatter{tf, t} }

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

func (tf timeFormatter) AppendSllm(buf []byte) []byte {
	if tf.fmt&TUTC != 0 {
		tf.t = tf.t.UTC()
	}

	if tf.fmt&(Tdate|Tyear|Tweekday) != 0 {
		ye, mo, dy := tf.t.Date()
		if tf.fmt&Tyear != 0 {
			buf = uitoa(buf, ye, 4)
			buf = append(buf, '-')
		}
		buf = uitoa(buf, int(mo), 2)
		buf = append(buf, '-')
		buf = uitoa(buf, dy, 2)
		if tf.fmt&Tweekday != 0 {
			buf = append(buf, ' ')
			buf = append(buf, tf.t.Weekday().String()[:2]...)
		}
		if tf.fmt&(Tclock|Tmillis|Tmicros) != 0 {
			buf = append(buf, ' ')
		}
	}

	if tf.fmt&(Tclock|Tmillis|Tmicros) != 0 {
		ho, mi, sc := tf.t.Clock()
		buf = uitoa(buf, ho, 2)
		buf = append(buf, ':')
		buf = uitoa(buf, mi, 2)
		buf = append(buf, ':')
		buf = uitoa(buf, sc, 2)
		if tf.fmt&Tmicros != 0 {
			buf = append(buf, '.')
			buf = uitoa(buf, tf.t.Nanosecond()/1000, 6)
		} else if tf.fmt&Tmillis != 0 {
			buf = append(buf, '.')
			buf = uitoa(buf, tf.t.Nanosecond()/1000000, 3)
		}
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
