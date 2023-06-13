package lumina

import (
	"context"
)

type heloHandler struct {
	serverSession *ServerSession
}

// Accept PKT_HELO only.
func (*heloHandler) AcceptRequest(t PacketType) bool {
	return t == PKT_HELO
}

// Translate PKT_HELO only.
func (*heloHandler) GetPacketOfType(t PacketType) Packet {
	switch t {
	case PKT_HELO:
		return &HeloPacket{}
	default:
		return nil
	}
}

// Serve HeloPacket only.
func (h *heloHandler) ServeRequest(ctx context.Context, req Request) (rsp Packet, err error) {
	helo := req.(*HeloPacket)
	logger := GetLogger(ctx)
	logger.Printf("protocol version: %d\nlicense: %v\n%s",
		helo.ClientVersion,
		helo.LicenseId,
		helo.Key,
	)
	h.serverSession.ctx = setProtocolVersion(ctx, helo.ClientVersion)
	if helo.ClientVersion < 5 {
		rsp = &RpcOkPacket{}
	} else {
		rsp = &HeloResultPacket{
			User: LuminaUser{
				LicenseInfo: UserLicenseInfo{
					Id:    "",
					Name:  "",
					Email: "",
				},
				Name:       "",
				Karma:      0,
				LastActive: UtcTimestamp(0),
				Features:   0,
			},
		}
	}
	return
}
