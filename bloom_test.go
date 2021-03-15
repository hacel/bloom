package bloom

import (
	"hash"
	"hash/fnv"
	"os"
	"reflect"
	"testing"

	"github.com/edsrzf/mmap-go"
	"github.com/spaolacci/murmur3"
)

func TestNew(t *testing.T) {
	type args struct {
		m uint32
	}
	tests := []struct {
		name string
		args args
		want *Filter
	}{
		{"t1", args{32}, &Filter{make([]byte, 4), []hash.Hash32{fnv.New32(), murmur3.New32()}, 32}},
		{"t2", args{0}, &Filter{make([]byte, 1), []hash.Hash32{fnv.New32(), murmur3.New32()}, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddCheck(t *testing.T) {
	type args struct {
		item []byte
	}
	tests := []struct {
		name string
		f    *Filter
		args args
	}{
		{"t1", New(32000), args{[]byte("message1")}},
		{"t2", New(1), args{[]byte("message2")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.Add(tt.args.item)
			if !tt.f.Check(tt.args.item) {
				t.Errorf("Check() = \"%s\" was not found in filter", tt.args.item)
			}
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	f := New(32768)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Add([]byte{byte(i)})
	}
}

func BenchmarkCheck(b *testing.B) {
	f := New(32768)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Check([]byte{byte(i)})
	}
}

func TestNewMMap(t *testing.T) {
	type args struct {
		m        uint32
		filename string
	}
	tests := []struct {
		name string
		args args
		want *MMapFilter
	}{
		{"t1", args{32, "mmap_t1.test"}, &MMapFilter{make(mmap.MMap, 4), []hash.Hash32{fnv.New32(), murmur3.New32()}, 32}},
		{"t2", args{0, "mmap_t2.test"}, &MMapFilter{make(mmap.MMap, 1), []hash.Hash32{fnv.New32(), murmur3.New32()}, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMMap(tt.args.m, tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMMap() = %v, want %v", got, tt.want)
			}
		})
		os.Remove(tt.args.filename)
	}
}

func TestMMapAddCheck(t *testing.T) {
	args := []struct {
		m        int
		filename string
	}{
		{32000, "mmap.test"},
	}
	for _, tt := range args {
		f := NewMMap(uint32(tt.m), tt.filename)
		for i := 0; i < tt.m; i++ {
			f.Add([]byte{byte(i)})
			if !f.Check([]byte{byte(i)}) {
				t.Errorf("Check() = \"%d\" was not found in filter", i)
			}
		}
		os.Remove(tt.filename)
	}
}

func BenchmarkMMapAdd(b *testing.B) {
	f := NewMMap(32768, "benchmmap.test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Add([]byte{byte(i)})
	}
	os.Remove("benchmap.test")
}

func BenchmarkMMapCheck(b *testing.B) {
	f := NewMMap(32768, "benchmmap.test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Check([]byte{byte(i)})
	}
	os.Remove("benchmap.test")
}
