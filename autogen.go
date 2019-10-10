//go:generate go run ./autogen

package lumina

import (
	"bytes"
	"io"
	"strings"
)

// This implementation is much more restrict than original one.
// This implementation does not support edge cases related to EOF.
func readUint32(r Reader) (v uint32, err error) {
	var b0, b1, b2, b3, b4 byte
	if b0, err = r.ReadByte(); err != nil {
		return
	}
	if b0&0x80 == 0 {
		// 0x00000000 ~ 0x0000007F
		v = uint32(b0)
		return
	}
	if b1, err = r.ReadByte(); err != nil {
		return
	}
	if b0&0x40 == 0 {
		// 0x00000080 ~ 0x00003FFF
		v = uint32(b0&0x3F)<<8 | uint32(b1)
		if v < 0x80 {
			err = errUnexpectedEncoding
		}
		return
	}
	if b2, err = r.ReadByte(); err != nil {
		return
	}
	if b3, err = r.ReadByte(); err != nil {
		return
	}
	if b0&0x20 == 0 {
		// 0x00004000 ~ 0x1FFFFFFF
		v = uint32(b0&0x1F)<<24 | uint32(b1)<<16 | uint32(b2)<<8 | uint32(b3)
		if v < 0x4000 {
			err = errUnexpectedEncoding
		}
		return
	}
	if b4, err = r.ReadByte(); err != nil {
		return
	}
	// 0x20000000 ~ 0xFFFFFFFF
	v = uint32(b1)<<24 | uint32(b2)<<16 | uint32(b3)<<8 | uint32(b4)
	if v < 0x20000000 /* || b0 & 0x1F != 0 */ {
		err = errUnexpectedEncoding
	}
	return
}

func writeUint32(w Writer, v uint32) (err error) {
	switch {
	case v <= 0x0000007F:
	case v <= 0x00003FFF:
		if err = w.WriteByte(0x80 | byte(v>>8)); err != nil {
			return
		}
	case v <= 0x1FFFFFFF:
		if err = w.WriteByte(0xC0 | byte(v>>24)); err != nil {
			return
		}
		if err = w.WriteByte(byte(v >> 16)); err != nil {
			return
		}
		if err = w.WriteByte(byte(v >> 8)); err != nil {
			return
		}
	default:
		if err = w.WriteByte(0xE0); err != nil {
			return
		}
		if err = w.WriteByte(byte(v >> 24)); err != nil {
			return
		}
		if err = w.WriteByte(byte(v >> 16)); err != nil {
			return
		}
		if err = w.WriteByte(byte(v >> 8)); err != nil {
			return
		}
	}
	err = w.WriteByte(byte(v))
	return
}

func readInt32(r Reader) (int32, error) {
	v, err := readUint32(r)
	return int32(v), err
}

func writeInt32(w Writer, v int32) error {
	return writeUint32(w, uint32(v))
}

func readUint64(r Reader) (uint64, error) {
	lo, err := readUint32(r)
	if err != nil {
		return 0, err
	}
	hi, err := readUint32(r)
	if err != nil {
		return 0, err
	}
	return uint64(hi)<<32 | uint64(lo), nil
}

func writeUint64(w Writer, v uint64) (err error) {
	if err = writeUint32(w, uint32(v)); err != nil {
		return
	}
	if err = writeUint32(w, uint32(v>>32)); err != nil {
		return
	}
	return
}

func readBool(r Reader) (bool, error) {
	t, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	if t > 1 {
		return false, errUnexpectedEncoding
	}
	return t != 0, nil
}

func writeBool(w Writer, v bool) error {
	b := byte(0)
	if v {
		b = 1
	}
	return w.WriteByte(b)
}

func readString(r Reader) (string, error) {
	buf := bytes.Buffer{}
	for {
		c, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if c == 0 {
			break
		}
		err = buf.WriteByte(c)
		if err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func writeString(w Writer, s string) error {
	if pos := strings.IndexByte(s, 0); pos != -1 {
		return errInvalidValue
	}
	if _, err := strings.NewReader(s).WriteTo(w); err != nil {
		return err
	}
	return w.WriteByte(0)
}

func readBytes(r Reader, s []byte) error {
	_, err := io.ReadFull(r, s)
	return err
}

func writeBytes(w Writer, s []byte) error {
	_, err := io.Copy(w, bytes.NewReader(s))
	return err
}
