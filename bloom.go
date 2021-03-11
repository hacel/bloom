package bloom

import (
	"hash"
	"hash/fnv"

	"github.com/spaolacci/murmur3"
)

// Filter is the representation of the bloom filter
type Filter struct {
	bitset  []byte
	hashers []hash.Hash32
	m       uint32 /* length of the filter in bits */
}

func max(x, y uint32) uint32 {
	if x < y {
		return y
	}
	return x
}

// New returns a new Filter
func New(m uint32) *Filter {
	m = max(1, m)
	return &Filter{
		make([]byte, max(1, 1+(m-1)/8)),
		[]hash.Hash32{fnv.New32(), murmur3.New32()},
		m,
	}
}

func (f *Filter) hash(item []byte) []uint32 {
	var ret []uint32
	for _, h := range f.hashers {
		h.Write(item)
		ret = append(ret, h.Sum32())
		h.Reset()
	}
	return ret
}

// Add adds a new item to the filter
func (f *Filter) Add(item []byte) {
	hashes := f.hash(item)
	for _, v := range hashes {
		pos := v % f.m
		idx, offset := (pos / 8), (pos % 8)
		f.bitset[idx] |= (1 << offset)
	}
}

// Check tests whether an item is in the filter
func (f *Filter) Check(item []byte) bool {
	exists := 0
	hashes := f.hash(item)
	for _, v := range hashes {
		pos := v % f.m
		idx, offset := (pos / 8), (pos % 8)
		if f.bitset[idx]&(1<<offset) != 0 {
			exists++
		}
	}
	return exists == len(hashes)
}
