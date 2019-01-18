//autogen RpcOkPacket PKT_RPC_OK
//autogen RpcFailPacket PKT_RPC_FAIL
//autogen RpcNotifyPacket PKT_RPC_NOTIFY

package lumina

type RpcOkPacket struct {
}

func (pkt *RpcOkPacket) validateFields() error {
    return nil
}

func (pkt *RpcOkPacket) getCache() []byte {
    return []byte{
        0x00, 0x00, 0x00, 0x00,
        byte(PKT_RPC_OK),
    }
}

type RpcFailPacket struct {
    packetCache
    Result      int32
    Error       string
}

func (pkt *RpcFailPacket) validateFields() error {
    return nil
}

type RpcNotifyPacket struct {
    packetCache
    Code        int32  // FIXME: or uint32?
    Message     string
}

func (pkt *RpcNotifyPacket) validateFields() error {
    return nil
}
