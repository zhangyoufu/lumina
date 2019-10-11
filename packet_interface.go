package lumina

type Packet interface {
	Serializable
	getType() PacketType
	validateFields() error
}
