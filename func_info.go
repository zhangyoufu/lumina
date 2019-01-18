//autogen FuncInfo
//autogen FuncInfoAndFrequency
//autogen FuncInfoAndPattern

package lumina

type FuncInfo struct {
    Name            string
    Size            uint32
    Metadata        []byte // *metadata_t: key, size, data
}

type FuncInfoAndFrequency struct {
    FuncInfo
    Frequency       uint32
}

type FuncInfoAndPattern struct {
    FuncInfo
    PatternId       PatternId
}
