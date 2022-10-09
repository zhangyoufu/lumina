//autogen HeloPacket PKT_HELO PKT_RPC_OK

package lumina

import "github.com/palantir/stacktrace"

// A HeloPacket is sent to server once after connection established.
type HeloPacket struct {
	packetCache
	ClientVersion int32 // int
	Key           LicenseKey
	LicenseId     LicenseId
	RecordConv    bool // unknown meaning
	Username      string
	Password      string
}

func (pkt *HeloPacket) validateFields() error {
	if pkt.ClientVersion < 2 || pkt.ClientVersion > 3 {
		return stacktrace.NewError("HeloPacket.ClientVersion=%v is unexpected",
			pkt.ClientVersion,
		)
	}
	if pkt.RecordConv {
		return stacktrace.NewError("HeloPacket.RecordConv=true is unexpected")
	}
	if pkt.ClientVersion < 3 {
		if pkt.Username != "" || pkt.Password != "" {
			return stacktrace.NewError("HeloPacket unexpected credential username=%q password=%q",
				pkt.Username, pkt.Password,
			)
		}
	}
	// TODO: check Key versus LicenseId
	return nil
}

func newHeloPacket(version int32, licKey LicenseKey, licId LicenseId) (pkt *HeloPacket) {
	pkt = &HeloPacket{
		ClientVersion: version,
		Key:           licKey,
		LicenseId:     licId,
		RecordConv:    false,
		Username:      "",
		Password:      "",
	}
	return
}
