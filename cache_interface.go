package lumina

// A mutable cache.
type cache interface {
    cacheReader
    cacheWriter
}

type cacheReader interface {
    // The return value could be nil, which means the cache is invalid.
    getCache() RawPacket
}

type cacheWriter interface {
    // data == nil means invalidate existing cache.
    //
    // SetCache takes ownership of data argument, caller should not change its
    // content after a successful call to SetCache.
    setCache(RawPacket)
}
