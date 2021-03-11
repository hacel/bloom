package bloom

import (
	"hash"
	"hash/fnv"
	"os"

	"github.com/edsrzf/mmap-go"
	"github.com/spaolacci/murmur3"
)

// MMapFilter implements a bloom filter on a memory map
type MMapFilter struct {
	Filter
	MMap mmap.MMap
}

// NewMMap returns a new MMapFilter
func NewMMap(m uint32, filename string) *MMapFilter {
	file, err := os.OpenFile("map.mmap", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil
	}
	mp, err := mmap.Map(file, mmap.RDWR, 0)
	if err != nil {
		return nil
	}
	m = max(1, m)
	return &MMapFilter{
		Filter{
			make([]byte, 1+(m-1)/8),
			[]hash.Hash32{fnv.New32(), murmur3.New32()},
			m,
		},
		mp,
	}
}
