//autogen PushMdResultPacket PKT_PUSH_MD_RESULT

package lumina

type PushMdResultPacket struct {
    packetCache
    Codes   []OpResult
}

func (*PushMdResultPacket) validateFields() error {
    return nil
}
