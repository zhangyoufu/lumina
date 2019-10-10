package lumina

import (
	"bytes"
	"encoding/binary"
	"io"
)

// A packet on-wire is composed of following parts:
//  4 bytes big-endian payload length
//  1 byte packet type
//  variable-length payload
type RawPacket []byte

func NewRawPacket(t PacketType, payload []byte) (pkt RawPacket, err error) {
	if len(payload) > MaximumLength {
		err = errInvalidPayloadLength
		return
	}

	length := int32(len(payload))
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, &length)
	buf.WriteByte(byte(t))
	buf.Write(payload)
	pkt = RawPacket(buf.Bytes())
	return
}

// Read one packet from the given io.Reader.
func (pkt *RawPacket) ReadFrom(r io.Reader) (err error) {
	buf := &bytes.Buffer{}
	if _, err = io.CopyN(buf, r, 5); err != nil {
		return
	}

	var length int32
	_ = binary.Read(bytes.NewReader(buf.Bytes()), binary.BigEndian, &length)
	if length < 0 || length > MaximumLength {
		err = errInvalidPayloadLength
		return
	}

	if _, err = io.CopyN(buf, r, int64(length)); err != nil {
		return
	}

	*pkt = RawPacket(buf.Bytes())
	return
}

func (pkt RawPacket) WriteTo(w io.Writer) (err error) {
	_, err = io.Copy(w, bytes.NewReader([]byte(pkt)))
	return
}

func (pkt RawPacket) GetType() PacketType {
	return PacketType(pkt[4])
}

// Return only the payload part as []byte.
func (pkt RawPacket) GetPayload() []byte {
	return pkt[5:]
}
