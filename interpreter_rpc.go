package lumina

type rpcInterpreterType struct{}

var rpcInterpreter = &rpcInterpreterType{}

// Interpret PKT_RPC_OK/PKT_RPC_FAIL/PKT_RPC_NOTIFY.
func (*rpcInterpreterType) GetPacketOfType(t PacketType) Packet {
	switch t {
	case PKT_RPC_OK:
		return &RpcOkPacket{}
	case PKT_RPC_FAIL:
		return &RpcFailPacket{}
	case PKT_RPC_NOTIFY:
		return &RpcNotifyPacket{}
	default:
		return nil
	}
}
