package lumina

var _ cache = &packetCache{}

type packetCache struct {
    cache   RawPacket
}

func (pkt *packetCache) getCache() RawPacket {
    return pkt.cache
}

func (pkt *packetCache) setCache(cache RawPacket) {
    pkt.cache = cache
}
