package lumina

import "time"

type QTime64 uint64

func (t QTime64) Seconds() uint32 {
	return uint32(t >> 32)
}

func (t QTime64) Microseconds() uint32 {
	return uint32(t)
}

func (t QTime64) GoTime() time.Time {
	return time.Unix(int64(t.Seconds()), int64(t.Microseconds())*1000)
}

func (t QTime64) String() string {
	return t.GoTime().Format(time.DateTime)
}

type UtcTimestamp QTime64

func (t UtcTimestamp) String() string {
	return QTime64(t).String() + " UTC"
}

func (t UtcTimestamp) GoString() string {
	return `"` + t.String() + `"`
}
