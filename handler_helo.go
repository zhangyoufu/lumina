package lumina

import (
	"context"
)

type heloHandlerType struct{}

var heloHandler = &heloHandlerType{}

// Accept PKT_HELO only.
func (*heloHandlerType) AcceptRequest(t PacketType) bool {
	return t == PKT_HELO
}

// Translate PKT_HELO only.
func (*heloHandlerType) GetPacketOfType(t PacketType) Packet {
	switch t {
	case PKT_HELO:
		return &HeloPacket{}
	default:
		return nil
	}
}

// Serve HeloPacket only.
func (*heloHandlerType) ServeRequest(ctx context.Context, req Request) (rsp Packet, err error) {
	helo := req.(*HeloPacket)
	GetLogger(ctx).Printf("license: %v\n%s",
		helo.LicenseId,
		helo.Key,
	)
	rsp = &RpcOkPacket{}
	return
}
