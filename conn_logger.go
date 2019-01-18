package lumina

import (
    "log"
    "net"
    "os"
)

func newConnLogger(conn net.Conn) *log.Logger {
    id := getConnId(conn)
    return log.New(os.Stderr, "["+id+"] ", log.LstdFlags)
}
