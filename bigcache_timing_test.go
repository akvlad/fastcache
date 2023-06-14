package fastcache

import (
	"sync"
	"testing"
)

const items = 100

func BenchmarkSetBig(b *testing.B) {
	value := createValue(8*1024*1024, 0)
	c := New(1024 * 1024)
	b.SetBytes(int64(len(value)))
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := uint64(0); i < items; i++ {
				c.SetBig(i, value)
			}
		}
	})
}

func BenchmarkGetBig(b *testing.B) {
	key := uint64(12345)
	value := createValue(8*1024*1024, 0)
	c := New(1024 * 1024)
	c.SetBig(key, value)
	b.SetBytes(int64(len(value)))
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := uint64(0); i < items; i++ {
				_ = c.GetBig(nil, 12345)
			}
		}
	})
}

func BenchmarkGetSetBig(b *testing.B) {
	key := uint64(12345)
	value := createValue(8*1024*1024, 0)
	c := New(1024 * 1024)
	c.SetBig(key, value)
	b.SetBytes(int64(len(value)))
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := uint64(0); i < items; i++ {
				c.SetBig(i, value)
			}
			for i := uint64(0); i < items; i++ {
				_ = c.GetBig(nil, i)
			}
		}
	})
}

func BenchmarkGetBigMap(b *testing.B) {
	key := uint64(12345)
	value := createValue(8*1024*1024, 0)
	c := make(map[uint64][]byte)
	c[key] = value
	mtx := &sync.Mutex{}
	b.SetBytes(int64(len(value)))
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := uint64(0); i < items; i++ {
				mtx.Lock()
				_ = c[key]
				mtx.Unlock()
			}
		}
	})

}

func BenchmarkSetBigMap(b *testing.B) {
	value := createValue(8*1024*1024, 0)
	mtx := sync.Mutex{}
	c := map[uint64][]byte{}
	b.SetBytes(int64(len(value)))
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := uint64(0); i < items; i++ {
				mtx.Lock()
				c[i] = value
				mtx.Unlock()
			}
		}
	})
}

func BenchmarkGetSetBigMap(b *testing.B) {
	value := createValue(8*1024*1024, 0)
	mtx := sync.Mutex{}
	c := map[uint64][]byte{}
	b.SetBytes(int64(len(value)))
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := uint64(0); i < items; i++ {
				mtx.Lock()
				c[i] = value
				mtx.Unlock()
			}
			for i := uint64(0); i < items; i++ {
				mtx.Lock()
				_ = c[i]
				mtx.Unlock()
			}
		}
	})
}
