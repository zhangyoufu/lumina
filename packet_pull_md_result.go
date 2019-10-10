//autogen PullMdResultPacket PKT_PULL_MD_RESULT

package lumina

type PullMdResultPacket struct {
	packetCache
	Codes   []OpResult
	Results []FuncInfoAndFrequency
}

func (pkt *PullMdResultPacket) validateFields() error {
	return nil
}
