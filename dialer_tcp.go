package lumina

import "net"

// TCPDialer embeds net.Dialer.
type TCPDialer struct {
	net.Dialer
	hideDialContext
	Addr string
}

// By embedding this struct, net.Dialer.DialContext is effectively hidden.
type hideDialContext struct {
	DialContext struct{}
}

// WARNING: Timeout is divided for each resolved address. Please refer to the
// documentation of net.(*Dialer).DialContext.
func (d *TCPDialer) Dial() (net.Conn, error) {
	return d.Dialer.Dial("tcp", d.Addr)
}
