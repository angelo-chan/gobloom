package gobloom

import (
	"github.com/spaolacci/murmur3"
)

// Strategy a strategy to translate byte array to k bit indexes.
type Strategy interface {
	// Indexes gets indexes.
	Indexes(data []byte, m uint, k uint) []uint
}

// MURMUR128MITZ64
type MURMUR128MITZ64 struct {
}

// Indexes gets indexes.
func (strategy *MURMUR128MITZ64) Indexes(data []byte, m uint, k uint) []uint {
	indexes := make([]uint, k)
	h1, h2 := murmur3.Sum128(data)
	combined := int64(h1)

	for i := uint(0); i < k; i++ {
		indexes[i] = uint((combined & int64(0x7fffffffffffffff)) % int64(m))
		combined += int64(h2)
	}
	return indexes
}
