// +build go1.13

package lumina

import (
	"crypto/tls"
	"crypto/x509"
)

// Clients can be reused instead of created as needed. Clients are safe for
// concurrent use by multiple goroutines.
type Client struct {
	Dialer     Dialer
	LicenseKey LicenseKey
	LicenseId  LicenseId
}

func (c *Client) getDialer() Dialer {
	if c.Dialer == nil {
		return defaultDialer
	} else {
		return c.Dialer
	}
}

const (
	hexRaysAddr = "lumina.hex-rays.com:443"
	hexRaysCert = `
-----BEGIN CERTIFICATE-----
MIIBtzCCAV2gAwIBAgIJAK3otIT/2KiZMAoGCCqGSM49BAMCMFMxCzAJBgNVBAYT
AkJFMQ8wDQYDVQQHDAZMacOoZ2UxFTATBgNVBAoMDEhleC1SYXlzIFNBLjEcMBoG
A1UEAwwTbHVtaW5hLmhleC1yYXlzLmNvbTAeFw0xODEwMTAxMjMwNTZaFw0xOTEw
MTAxMjMwNTZaMFMxCzAJBgNVBAYTAkJFMQ8wDQYDVQQHDAZMacOoZ2UxFTATBgNV
BAoMDEhleC1SYXlzIFNBLjEcMBoGA1UEAwwTbHVtaW5hLmhleC1yYXlzLmNvbTBZ
MBMGByqGSM49AgEGCCqGSM49AwEHA0IABG2TLxpVsgDiji3F5OlYJQblgj8jYHTV
WkpFxJdIKQXXEWt0AmSJEHrtWGgnQWPL+8vhu2JY875lsfkuFaoOlxajGjAYMBYG
A1UdJQEB/wQMMAoGCCsGAQUFBwMBMAoGCCqGSM49BAMCA0gAMEUCIQCi2eEmOpbA
339TMowwLvT3KN8R+C/55tY6CPwIWSWhIwIgGsZ4k15nl7RSUydPQHhhT37W119W
Tp4/G0S1wzuXibA=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIBwTCCAWigAwIBAgIUTywOBIR2odB59aEjU981FBmOi+AwCgYIKoZIzj0EAwIw
UzELMAkGA1UEBhMCQkUxDzANBgNVBAcMBkxpw6hnZTEVMBMGA1UECgwMSGV4LVJh
eXMgU0EuMRwwGgYDVQQDDBNsdW1pbmEuaGV4LXJheXMuY29tMB4XDTE5MTAwODE0
MTg1OFoXDTIwMTAwNzE0MTg1OFowUzELMAkGA1UEBhMCQkUxDzANBgNVBAcMBkxp
w6hnZTEVMBMGA1UECgwMSGV4LVJheXMgU0EuMRwwGgYDVQQDDBNsdW1pbmEuaGV4
LXJheXMuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEbZMvGlWyAOKOLcXk
6VglBuWCPyNgdNVaSkXEl0gpBdcRa3QCZIkQeu1YaCdBY8v7y+G7YljzvmWx+S4V
qg6XFqMaMBgwFgYDVR0lAQH/BAwwCgYIKwYBBQUHAwEwCgYIKoZIzj0EAwIDRwAw
RAIgB6B+bFSXowi5wV0xJXsCyyR/EjKg1OIHlFbDW9SHCRoCIH+b7xguFt0IptGV
qx1spjBjuLXas8sMFJKDqheggBl3
-----END CERTIFICATE-----
`
)

var defaultDialer Dialer

func init() {
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM([]byte(hexRaysCert)); !ok {
		panic("unable to parse Hex-Rays cert")
	}

	d := &TLSDialer{}
	d.Addr = hexRaysAddr
	d.RootCAs = roots
	d.MinVersion = tls.VersionTLS13

	defaultDialer = d
}
