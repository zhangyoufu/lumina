package lumina

import "encoding/hex"

// md5_t
type MD5Digest [16]byte

func (d MD5Digest) String() string {
	return hex.EncodeToString(d[:])
}
