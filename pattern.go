//autogen PatternId

package lumina

type PatternType uint32
const (
    PAT_TYPE_UNKNOWN PatternType = 0
    PAT_TYPE_MD5     PatternType = 1
)

type PatternId struct {
    Type            PatternType
    Data            []byte
}
