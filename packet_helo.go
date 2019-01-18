//autogen HeloPacket PKT_HELO PKT_RPC_OK

package lumina

import "github.com/palantir/stacktrace"

// A HeloPacket is sent to server once after connection established.
type HeloPacket struct {
    packetCache
    ClientVersion   int32       // int
    Key             LicenseKey
    LicenseId       LicenseId
    RecordConv      bool        // unknown meaning
}

func (pkt *HeloPacket) validateFields() error {
    if pkt.ClientVersion != 2 {
        return stacktrace.NewError("HeloPacket.ClientVersion=%v is unexpected",
            pkt.ClientVersion,
        )
    }
    if pkt.RecordConv {
        return stacktrace.NewError("HeloPacket.RecordConv=true is unexpected")
    }
    // TODO: check Key versus LicenseId
    return nil
}

func newHeloPacket(licKey LicenseKey, licId LicenseId) (pkt *HeloPacket) {
    pkt = &HeloPacket{
        ClientVersion:  2,
        Key:            licKey,
        LicenseId:      licId,
        RecordConv:     false,
    }
    return
}
