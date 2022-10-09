//autogen DecompilePacket PKT_DECOMPILE PKT_DECOMPILE_RESULT

package lumina

type DecompilePacket struct {
	packetCache
	Opaque []byte
}

func (pkt *DecompilePacket) validateFields() error {
	// TODO
	return nil
}
