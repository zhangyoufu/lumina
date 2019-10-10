package lumina

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

type LicenseId [6]byte

// TODO: more strict
func ParseLicenseId(s string) (id LicenseId) {
	data, _ := hex.DecodeString(strings.Replace(s, "-", "", 3))
	copy(id[:], data)
	return
}

func (id LicenseId) String() string {
	return fmt.Sprintf("%02X-%02X%02X-%02X%02X-%02X",
		id[0], id[1], id[2], id[3], id[4], id[5],
	)
}

type LicenseKey []byte

func (lic LicenseKey) String() string {
	return string(lic)
}

type IDALicenseInfo struct {
	Id LicenseId
	OS byte
}

func (lic LicenseKey) GetIDAInfo() *IDALicenseInfo {
	r := bufio.NewReader(bytes.NewReader(lic))
	var err error
	for err == nil {
		var line string
		line, err = r.ReadString('\n')
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		licenseId := fields[0]
		product := fields[1]
		if !strings.HasPrefix(product, "IDA") {
			continue
		}
		return &IDALicenseInfo{
			Id: ParseLicenseId(licenseId),
			OS: product[len(product)-1],
		}
	}
	return nil
}
