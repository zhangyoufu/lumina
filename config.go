package lumina

import "math"

// This constant limits the maximum length of a packet payload.
//
// Currently this constant is set to math.MaxInt32. But notice that there is a
// 4096-bytes limitation implemented whilst not enabled in IDA Pro 7.2.181105.
// So this constant may be changed in future.
const MaximumLength = math.MaxInt32
