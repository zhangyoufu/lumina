//autogen HeloResultPacket PKT_HELO_RESULT
//autogen LuminaUser
//autogen UserLicenseInfo

package lumina

type UserLicenseInfo struct {
	Id    string
	Name  string
	Email string
}

type LuminaUser struct {
	LicenseInfo UserLicenseInfo
	Name        string
	Karma       int32
	LastActive  UtcTimestamp
	Features    uint32 // is_admin, can_del_history
}

// A HeloResultPacket is sent from server as response to HeloPacket. Only
// available for client version ≥ 5.
type HeloResultPacket struct {
	packetCache
	User LuminaUser
}

func (pkt *HeloResultPacket) validateFields() error {
	// TODO
	return nil
}
