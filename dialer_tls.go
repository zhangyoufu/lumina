package lumina

import (
    "crypto/tls"
    "net"
)

// TLSDialer embeds net.Dialer and tls.Config.
type TLSDialer struct {
    TCPDialer
    tls.Config
}

// If RootCAs is nil, the host's root CA set will be used.
func (d *TLSDialer) Dial() (net.Conn, error) {
    return tls.DialWithDialer(&d.Dialer, "tcp", d.Addr, &d.Config)
}
