package fastcache

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/allegro/bigcache"
)

func BenchmarkBigCacheSet(b *testing.B) {
	const items = 1 << 16
	cfg := bigcache.DefaultConfig(time.Minute)
	cfg.Verbose = false
	c, err := bigcache.NewBigCache(cfg)
	if err != nil {
		b.Fatalf("cannot create cache: %s", err)
	}
	defer c.Close()
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xyza")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				if err := c.Set(b2s(k), v); err != nil {
					panic(fmt.Errorf("unexpected error: %s", err))
				}
			}
		}
	})
}

func BenchmarkBigCacheGet(b *testing.B) {
	const items = 1 << 16
	cfg := bigcache.DefaultConfig(time.Minute)
	cfg.Verbose = false
	c, err := bigcache.NewBigCache(cfg)
	if err != nil {
		b.Fatalf("cannot create cache: %s", err)
	}
	defer c.Close()
	k := []byte("\x00\x00\x00\x00")
	v := []byte("xyza")
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		if err := c.Set(b2s(k), v); err != nil {
			b.Fatalf("unexpected error: %s", err)
		}
	}

	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				vv, err := c.Get(b2s(k))
				if err != nil {
					panic(fmt.Errorf("BUG: unexpected error: %s", err))
				}
				if string(vv) != string(v) {
					panic(fmt.Errorf("BUG: invalid value obtained; got %q; want %q", vv, v))
				}
			}
		}
	})
}

func BenchmarkBigCacheSetGet(b *testing.B) {
	const items = 1 << 16
	cfg := bigcache.DefaultConfig(time.Minute)
	cfg.Verbose = false
	c, err := bigcache.NewBigCache(cfg)
	if err != nil {
		b.Fatalf("cannot create cache: %s", err)
	}
	defer c.Close()
	b.ReportAllocs()
	b.SetBytes(2 * items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xyza")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				if err := c.Set(b2s(k), v); err != nil {
					panic(fmt.Errorf("unexpected error: %s", err))
				}
			}
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				vv, err := c.Get(b2s(k))
				if err != nil {
					panic(fmt.Errorf("BUG: unexpected error: %s", err))
				}
				if string(vv) != string(v) {
					panic(fmt.Errorf("BUG: invalid value obtained; got %q; want %q", vv, v))
				}
			}
		}
	})
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func BenchmarkCacheSet(b *testing.B) {
	const items = 1 << 16
	c := New(12 * items)
	defer c.Reset()
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := uint64(1)
		v := []byte("xyza")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k++
				c.Set(k, v)
			}
		}
	})
}

func BenchmarkCacheGet(b *testing.B) {
	const items = 1 << 16
	c := New(12 * items)
	defer c.Reset()
	k := uint64(1)
	v := []byte("xyza")
	for i := 0; i < items; i++ {
		k++

		c.Set(k, v)
	}

	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		var buf []byte
		for pb.Next() {
			k := uint64(1)
			for i := 0; i < items; i++ {
				k++
				buf = c.Get(buf[:0], k)
				if string(buf) != string(v) {
					panic(fmt.Errorf("BUG: invalid value obtained; got %q; want %q", buf, v))
				}
			}
		}
	})
}

func BenchmarkCacheHas(b *testing.B) {
	const items = 1 << 16
	c := New(12 * items)
	defer c.Reset()
	k := uint64(1)
	for i := 0; i < items; i++ {
		k++
		c.Set(k, nil)
	}

	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := uint64(1)
			for i := 0; i < items; i++ {
				k++
				if !c.Has(k) {
					panic(fmt.Errorf("BUG: missing value for key %q", k))
				}
			}
		}
	})
}

func BenchmarkCacheSetGet(b *testing.B) {
	const items = 1 << 16
	c := New(12 * items)
	defer c.Reset()
	b.ReportAllocs()
	b.SetBytes(2 * items)
	b.RunParallel(func(pb *testing.PB) {

		v := []byte("xyza")
		var buf []byte
		for pb.Next() {
			k := uint64(1)
			for i := 0; i < items; i++ {
				k++
				c.Set(k, v)
			}
			k = 1
			for i := 0; i < items; i++ {
				k++
				buf = c.Get(buf[:0], k)
				if string(buf) != string(v) {
					panic(fmt.Errorf("BUG: invalid value obtained; got %q; want %q", buf, v))
				}
			}
		}
	})
}

func BenchmarkStdMapSet(b *testing.B) {
	const items = 1 << 16
	m := make(map[uint64][]byte)
	var mu sync.Mutex
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		v := []byte("xyza")
		for pb.Next() {
			k := uint64(1)
			for i := 0; i < items; i++ {
				k++
				mu.Lock()
				m[k] = v
				mu.Unlock()
			}
		}
	})
}

func BenchmarkStdMapGet(b *testing.B) {
	const items = 1 << 16
	m := make(map[uint64][]byte)
	k := uint64(1)
	v := []byte("xyza")
	for i := 0; i < items; i++ {
		k++
		m[k] = v
	}

	var mu sync.RWMutex
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := uint64(1)
			for i := 0; i < items; i++ {
				k++
				mu.RLock()
				vv := m[k]
				mu.RUnlock()
				if string(vv) != string(v) {
					panic(fmt.Errorf("BUG: unexpected value; got %q; want %q", vv, v))
				}
			}
		}
	})
}

func BenchmarkStdMapSetGet(b *testing.B) {
	const items = 1 << 16
	m := make(map[uint64][]byte)
	var mu sync.RWMutex
	b.ReportAllocs()
	b.SetBytes(2 * items)
	b.RunParallel(func(pb *testing.PB) {
		v := []byte("xyza")
		for pb.Next() {
			k := uint64(1)
			for i := 0; i < items; i++ {
				k++
				mu.Lock()
				m[k] = v
				mu.Unlock()
			}
			k = 1
			for i := 0; i < items; i++ {
				k++
				mu.RLock()
				vv := m[k]
				mu.RUnlock()
				if string(vv) != string(v) {
					panic(fmt.Errorf("BUG: unexpected value; got %q; want %q", vv, v))
				}
			}
		}
	})
}

func BenchmarkSyncMapSet(b *testing.B) {
	const items = 1 << 16
	m := sync.Map{}
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		v := "xyza"
		for pb.Next() {
			k := uint64(1)
			for i := 0; i < items; i++ {
				k++
				m.Store(k, v)
			}
		}
	})
}

func BenchmarkSyncMapGet(b *testing.B) {
	const items = 1 << 16
	m := sync.Map{}
	v := "xyza"
	k := uint64(1)
	for i := 0; i < items; i++ {
		k++
		m.Store(k, v)
	}

	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := uint64(1)
			for i := 0; i < items; i++ {
				k++
				vv, ok := m.Load(k)
				if !ok || vv.(string) != string(v) {
					panic(fmt.Errorf("BUG: unexpected value; got %q; want %q", vv, v))
				}
			}
		}
	})
}

func BenchmarkSyncMapSetGet(b *testing.B) {
	const items = 1 << 16
	m := sync.Map{}
	b.ReportAllocs()
	b.SetBytes(2 * items)
	b.RunParallel(func(pb *testing.PB) {
		v := "xyza"
		for pb.Next() {
			k := uint64(1)
			for i := 0; i < items; i++ {
				k++
				m.Store(k, v)
			}
			k = 1
			for i := 0; i < items; i++ {
				k++
				vv, ok := m.Load(k)
				if !ok || vv.(string) != string(v) {
					panic(fmt.Errorf("BUG: unexpected value; got %q; want %q", vv, v))
				}
			}
		}
	})
}
