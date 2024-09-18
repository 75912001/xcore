package util

import "hash/fnv"

// HASH32
func HASH32(data []byte) uint32 {
	h := fnv.New32()
	_, _ = h.Write(data)
	return h.Sum32()
}

// HASH64
func HASH64(data []byte) uint64 {
	h := fnv.New64()
	_, _ = h.Write(data)
	return h.Sum64()
}
