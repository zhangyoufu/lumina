//autogen PushMdPacket PKT_PUSH_MD PKT_PUSH_MD_RESULT
//autogen InputFile

package lumina

type PushMdFlag uint32
const (
    PMF_PUSH_MODE_MASK          PushMdFlag = 0xF
    PMF_PUSH_OVERRIDE_IF_BETTER PushMdFlag = 0x0
    PMF_PUSH_OVERRIDE           PushMdFlag = 0x1
    PMF_PUSH_DO_NOT_OVERRIDE    PushMdFlag = 0x2
    PMF_PUSH_MERGE              PushMdFlag = 0x3
)
func (flags PushMdFlag) GetMode() PushMdFlag {
    return flags & PMF_PUSH_MODE_MASK
}

type InputFile struct {
    Path            string
    MD5             MD5Digest
}

type PushMdPacket struct {
    packetCache
    Flags           PushMdFlag
    Idb             string
    Input           InputFile
    Hostname        string
    Contents        []FuncInfoAndPattern
    EAs             []uint64    // append_ea64 / unpack_ea64
}

func (*PushMdPacket) validateFields() error {
    return nil
}
