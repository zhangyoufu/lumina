package lumina

// Translate PacketType to Packet for unmarshalling incoming packets.
type Interpreter interface {
	// Returns nil if the given PacketType is not supported.
	GetPacketOfType(PacketType) Packet
}
