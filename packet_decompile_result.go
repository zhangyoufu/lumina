//autogen DecompileResultPacket PKT_DECOMPILE_RESULT

package lumina

type DecompileResultPacket struct {
	packetCache
	Payload  []byte
}

func (pkt *DecompileResultPacket) validateFields() error {
	return nil
}
