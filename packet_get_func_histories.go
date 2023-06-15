//autogen GetFuncHistoriesPacket PKT_GET_FUNC_HISTORIES PKT_GET_FUNC_HISTORIES_RESULT

package lumina

type GetFuncHistoriesPacket struct {
	packetCache
	PatternIds []PatternId
	Unknown    uint32
}

func (pkt *GetFuncHistoriesPacket) validateFields() error {
	// TODO
	return nil
}
