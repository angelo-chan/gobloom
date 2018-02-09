package gobloom

import (
	"encoding/json"
	"math"

	"github.com/willf/bitset"
)

// BloomFilter class wrapper for bloom filter.
// Cheat sheet:
//
// m: total bits
// n: expected insertions
// b: m/n, bits per insertion
// p: expected false positive probability
// k: number of hash functions
//
// 1) Optimal k = b * ln2
// 2) p = (1 - e ^ (-kn/m))^k
// 3) For optimal k: p = 2 ^ (-k) ~= 0.6185^b
// 4) For optimal k: m = -nlnp / ((ln2) ^ 2)
type BloomFilter struct {
	m    uint
	k    uint
	bits *bitset.BitSet
	s    Strategy
}

func pad(m uint) uint {
	r := m % 64
	d := m / 64
	if r == 0 {
		return m
	}
	return 64 * (d + 1)
}

// New creates a new bloom filter.
// n: expected insertions
// p: expected false positive probability
func New(n uint, p float64) *BloomFilter {
	m, k := EstimateParameters(n, p)
	return &BloomFilter{pad(m), k, bitset.New(m), &MURMUR128MITZ64{}}
}

// NewWithStrategy creates a new bloom filter with strategy.
// n: expected insertions
// p: expected false positive probability
// s: index strategy
func NewWithStrategy(n uint, p float64, s Strategy) *BloomFilter {
	m, k := EstimateParameters(n, p)
	return &BloomFilter{pad(m), k, bitset.New(m), s}
}

// From populates a bloom filter.
// bits: long array represents bits
// k: number of hash functions
func From(bits []uint64, k uint) *BloomFilter {
	m := uint(len(bits) * 64)
	return &BloomFilter{m, k, bitset.From(bits), &MURMUR128MITZ64{}}
}

// FromWithStrategy populates a bloom filter with strategy.
// bits: long array represents bits
// k: number of hash functions
// s: index strategy
func FromWithStrategy(bits []uint64, k uint, s Strategy) *BloomFilter {
	m := uint(len(bits) * 64)
	return &BloomFilter{m, k, bitset.From(bits), s}
}

// EstimateParameters estimates requirements for m and k.
// n: expected insertions
// p: expected false positive probability
func EstimateParameters(n uint, p float64) (m uint, k uint) {
	m = uint(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2))
	k = uint(math.Max(float64(1), math.Floor(0.5+float64(m)/float64(n)*math.Log(2))))
	return
}

// M returns the capacity of a Bloom filter.
func (f *BloomFilter) M() uint {
	return f.m
}

// K returns the number of hash functions used in the BloomFilter.
func (f *BloomFilter) K() uint {
	return f.k
}

// Put data to the Bloom Filter. Returns the filter (allows chaining)
func (f *BloomFilter) Put(data []byte) *BloomFilter {
	indexes := f.s.Indexes(data, f.m, f.k)
	for _, index := range indexes {
		f.bits.Set(index)
	}
	return f
}

// MightContain query data existence.
func (f *BloomFilter) MightContain(data []byte) bool {
	indexes := f.s.Indexes(data, f.m, f.k)
	for _, index := range indexes {
		if !f.bits.Test(index) {
			return false
		}
	}
	return true
}

// BloomFilterJSON class wrapper for marshaling/unmarshaling BloomFilter struct.
type BloomFilterJSON struct {
	M uint     `json:"m"`
	K uint     `json:"k"`
	B []uint64 `json:"bits"`
}

// MarshalJSON implements json.Marshaler interface.
func (f *BloomFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(BloomFilterJSON{f.m, f.k, f.bits.Bytes()})
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (f *BloomFilter) UnmarshalJSON(data []byte) error {
	var fJSON BloomFilterJSON
	err := json.Unmarshal(data, &fJSON)
	if err != nil {
		return err
	}
	f.m = fJSON.M
	f.k = fJSON.K
	f.bits = bitset.From(fJSON.B)
	f.s = &MURMUR128MITZ64{}
	return nil
}

// SetStrategy sets strategy. Do not set it manually except only you unmarshal a bloom filter that using different strategy.
func (f *BloomFilter) SetStrategy(s Strategy) {
	f.s = s
}
