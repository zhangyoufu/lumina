package lumina

import "net"

// A Dialer dials like a lot.
type Dialer interface {
    Dial() (conn net.Conn, err error)
}
