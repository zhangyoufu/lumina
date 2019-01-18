package lumina

type rpcInterpreterType struct {}
var rpcInterpreter = &rpcInterpreterType{}

// Translate PKT_HELO only.
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
