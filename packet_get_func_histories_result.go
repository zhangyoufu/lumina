//autogen GetFuncHistoriesResultPacket PKT_GET_FUNC_HISTORIES_RESULT
//autogen FuncHistoryData
//autogen FuncHistoryBase
//autogen FuncHistory

package lumina

type FuncHistoryData struct {
	Name     string
	Metadata []byte
}

type FuncHistoryBase struct {
	Unknown1 uint64
	Unknown2 uint64
	Data     FuncHistoryData
}

type FuncHistory struct {
	FuncHistoryBase
	Timestamp  UtcTimestamp
	AuthorIdx  uint32 // guess
	IdbPathIdx uint32 // guess
}

type FuncHistories []FuncHistory

type GetFuncHistoriesResultPacket struct {
	packetCache
	Codes    []OpResult // guess
	Funcs    []FuncHistories
	Authors  []string // guess
	IdbPaths []string // guess
}

func (pkt *GetFuncHistoriesResultPacket) validateFields() error {
	// TODO
	return nil
}
