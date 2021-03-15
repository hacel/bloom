package bloom

import (
	"errors"
	"hash"
	"hash/fnv"
	"os"

	"github.com/edsrzf/mmap-go"
	"github.com/spaolacci/murmur3"
)

func max(x, y uint32) uint32 {
	if x < y {
		return y
	}
	return x
}

// Filter is the representation of the bloom filter
type Filter struct {
	Bitset  []byte
	hashers []hash.Hash32
	m       uint32 /* length of the filter in bits */
}

// New returns a new Filter
func New(m uint32) *Filter {
	m = max(1, m)
	return &Filter{
		make([]byte, 1+(m-1)/8),
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
		f.Bitset[idx] |= (1 << offset)
	}
}

// Check tests whether an item is in the filter
func (f *Filter) Check(item []byte) bool {
	exists := 0
	hashes := f.hash(item)
	for _, v := range hashes {
		pos := v % f.m
		idx, offset := (pos / 8), (pos % 8)
		if f.Bitset[idx]&(1<<offset) != 0 {
			exists++
		}
	}
	return exists == len(hashes)
}

// MMapFilter implements a bloom filter on a memory map
type MMapFilter struct {
	Bitset  mmap.MMap
	hashers []hash.Hash32
	m       uint32 /* length of the filter in bits */
}

func initFile(filename string, size uint32) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.FileMode(0o0644))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(make([]byte, size))
	if err != nil {
		return err
	}
	return nil
}

// NewMMap returns a new MMapFilter
func NewMMap(m uint32, filename string) *MMapFilter {
	m = max(1, m)
	mbyte := 1 + (m-1)/8

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		if err := initFile(filename, mbyte); err != nil {
			return nil
		}
	}

	file, err := os.OpenFile(filename, os.O_RDWR, os.FileMode(0o0644))
	if err != nil {
		return nil
	}
	defer file.Close()

	memMap, err := mmap.Map(file, mmap.RDWR, 0)
	if err != nil {
		return nil
	}
	defer memMap.Flush()

	return &MMapFilter{
		memMap,
		[]hash.Hash32{fnv.New32(), murmur3.New32()},
		m,
	}
}

func (m *MMapFilter) hash(item []byte) []uint32 {
	var ret []uint32
	for _, h := range m.hashers {
		h.Write(item)
		ret = append(ret, h.Sum32())
		h.Reset()
	}
	return ret
}

// Add adds a new item to the mmap filter
func (m *MMapFilter) Add(item []byte) {
	hashes := m.hash(item)
	for _, v := range hashes {
		pos := v % m.m
		idx, offset := (pos / 8), (pos % 8)
		m.Bitset[idx] |= (1 << offset)
	}
}

// Check tests whether an item is in the mmap filter
func (m *MMapFilter) Check(item []byte) bool {
	exists := 0
	hashes := m.hash(item)
	for _, v := range hashes {
		pos := v % m.m
		idx, offset := (pos / 8), (pos % 8)
		if m.Bitset[idx]&(1<<offset) != 0 {
			exists++
		}
	}
	return exists == len(hashes)
}
