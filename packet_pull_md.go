//autogen PullMdPacket PKT_PULL_MD PKT_PULL_MD_RESULT

package lumina

type PullMdPacket struct {
	packetCache
	Flags      MdKeyFlag
	Keys       []MdKey
	PatternIds []PatternId
}

func (pkt *PullMdPacket) validateFields() error {
	// TODO
	return nil
}
