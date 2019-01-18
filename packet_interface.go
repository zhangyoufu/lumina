package lumina

type Packet interface {
    Serializable
    getType() PacketType
    validateFields() error
}

type Request interface {
    Packet
    getResponseType() PacketType
}
