package lumina

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

func newTaggedLogger() *log.Logger {
	id := fmt.Sprintf("%08x", rand.Uint32())
	return log.New(os.Stderr, "["+id+"] ", log.LstdFlags)
}
