//go:build go1.13

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
MIIF0TCCA7mgAwIBAgIULzKtEOP9Q7V/L/G4Rnv4L3vq/hEwDQYJKoZIhvcNAQEN
BQAwVDELMAkGA1UEBhMCQkUxDzANBgNVBAcMBkxpw6hnZTEVMBMGA1UECgwMSGV4
LVJheXMgU0EuMR0wGwYDVQQDDBRIZXgtUmF5cyBTQS4gUm9vdCBDQTAeFw0yMDA1
MDQxMTAyMDhaFw00MDA0MjkxMTAyMDhaMFQxCzAJBgNVBAYTAkJFMQ8wDQYDVQQH
DAZMacOoZ2UxFTATBgNVBAoMDEhleC1SYXlzIFNBLjEdMBsGA1UEAwwUSGV4LVJh
eXMgU0EuIFJvb3QgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDB
rsEh48VNyjCPROSYzw5viAcfDuBoDAHe3bIRYMaGm2a6omSXSzT02RAipSlO6nJZ
/PgNEipaXYLbEXmrrGdnSdBu8ub51t17AdGcGYzzPjSIpIVH5mX2iObHdS3gyNzp
JKJQUCDM6FdJa8ZcztKw+bXsN1ftKaZCzHcuUBc8P5lkiRGcuYfbiHri5C02pGo1
3y4Oz99Sot8KUfwNhByOOGOweYyfn9NgmhqhkBu27+6rxpmuR7mHyOhfnLs+psQ0
yjE6bzul2ilWFrOSaLAxKbhBLLQDWCYeBvXmE0IzmZVbo2DqTU+NWREU6avmRRBz
6RnZHFUhl2LVbJ5Ar45BawR38bRNro6VNCTq89rBXVFeCnk9Ja6v4ZAoWmjJupHC
pXTIxoebkoeWAwICuz63cWsRh1y2aqdgQ6v9yVErA64GhgCkpJO82HDtA9Siqge3
T+rgUnj1pcllGKgxAFYcKhlCLl4+bm0ohlxF0WF8VMhG/TBLNH3MlJFjlMoBwQnl
APheEgZWoQSEjAkzRLUrRw7kVk/Qt8G5hFGLb3UjE8SKDPKRYSBAUN/uP8YHKFqo
2arpTCi1DO4SqX8r6zqzslVTf6uWTiq8MNkZ/+7NYr1/JPT25iMlw6sa6g4GUPpQ
zhRaPy19obGe43u4vjpyse9g5vqX9p3u9MI14x3k6QIDAQABo4GaMIGXMB0GA1Ud
DgQWBBQaxNacfM7XKjKIutIHrc6tjiE9DTAfBgNVHSMEGDAWgBQaxNacfM7XKjKI
utIHrc6tjiE9DTAPBgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQEAwIBhjA0BgNV
HR8ELTArMCmgJ6AlhiNodHRwOi8vY3JsLmhleC1yYXlzLmNvbS9yb290X2NhLmNy
bDANBgkqhkiG9w0BAQ0FAAOCAgEAdKp4InpRk5z0BjPs6CcJSgKbCH0MXZqbt/EM
/4dJPvmA6tAexJpv9e9BmT/DOB84QB2xzQlEiNOB7/V4j3oDij5mMwRyqYL24l3g
HAavwc+dLrpzX/54uZmH9bKs7yj3fk/vU3e7th720ArL2/YZjHV2Wx0BMcs+YVit
phvG2mxu16DTpidms3pCj25eEISJvXfe8XEfKOP1FxGCpmKxx6qPHlNASOp5zdwV
iEimkguUwzCsmmPI5rEWLXdLRxc0CkffmbsNmsF8SZz38CiwuRlichDDdZuJXji7
jnZF7h04Mo2AKPt6wJ9+66rYqDigvP9sHGKpQp5hr1DMukFGnei3S9h5Kp8eDhRX
Y24y/CJVNO0rxYoFPUnOwbSUF3Fwu4fX3Ezq5eW7N0Nl7s0XHExb/P9fmhPxQBV1
gwr665inq5ZwD8H9uwGEVp3IBT9cHRu8ieZrQDMI1UqPOy+2EWNPtY4KxmgerTbc
N0VH4BuE8tdxTGUckg4JTbsNRUbqxSXmSL9jA1dLBT63lbMLIU06dIdqNbpxE4GV
MgOLwqwx/BF+FZgQTttdjmpexml6NIDVGDBxfyECJ5vdwxbKMIRfo7fp0jRpjZpP
8bw4BPnx0Y4NpMzKxiWS0i7re9iEafdh6GtpNynKU0JFSKrIwmIecKF+Z4ZUE/1K
+t/FOgI=
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
