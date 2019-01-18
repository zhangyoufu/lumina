package lumina

import (
    "hash/fnv"
    "net"
    "strconv"
)

func getConnId(conn net.Conn) string {
    local := conn.LocalAddr()
    remote := conn.RemoteAddr()
    h := fnv.New32a()
    h.Write([]byte(remote.Network()))
    h.Write([]byte(remote.String()))
    h.Write([]byte(local.String()))
    return strconv.FormatUint(uint64(h.Sum32()), 16)
}
