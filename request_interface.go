package lumina

type Request interface {
	Packet
	getResponseType() PacketType
}
