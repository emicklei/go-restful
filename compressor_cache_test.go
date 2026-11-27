package restful

import (
	"compress/gzip"
	"compress/zlib"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

type boundedCompressorCase struct {
	name    string
	acquire func(*BoundedCachedCompressors) interface{}
	release func(*BoundedCachedCompressors, interface{})
	drain   func(*BoundedCachedCompressors) bool
}

func boundedCompressorCases() []boundedCompressorCase {
	return []boundedCompressorCase{
		{
			name:    "gzip-writer",
			acquire: func(b *BoundedCachedCompressors) interface{} { return b.AcquireGzipWriter() },
			release: func(b *BoundedCachedCompressors, v interface{}) { b.ReleaseGzipWriter(v.(*gzip.Writer)) },
			drain: func(b *BoundedCachedCompressors) bool {
				select {
				case <-b.gzipWriters:
					return true
				default:
					return false
				}
			},
		},
		{
			name:    "gzip-reader",
			acquire: func(b *BoundedCachedCompressors) interface{} { return b.AcquireGzipReader() },
			release: func(b *BoundedCachedCompressors, v interface{}) { b.ReleaseGzipReader(v.(*gzip.Reader)) },
			drain: func(b *BoundedCachedCompressors) bool {
				select {
				case <-b.gzipReaders:
					return true
				default:
					return false
				}
			},
		},
		{
			name:    "zlib-writer",
			acquire: func(b *BoundedCachedCompressors) interface{} { return b.AcquireZlibWriter() },
			release: func(b *BoundedCachedCompressors, v interface{}) { b.ReleaseZlibWriter(v.(*zlib.Writer)) },
			drain: func(b *BoundedCachedCompressors) bool {
				select {
				case <-b.zlibWriters:
					return true
				default:
					return false
				}
			},
		},
	}
}

func TestBoundedCachedCompressorsConcurrentRelease(t *testing.T) {
	old := runtime.GOMAXPROCS(8)
	defer runtime.GOMAXPROCS(old)

	for _, tc := range boundedCompressorCases() {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBoundedCachedCompressors(1, 1)
			// Acquire more resources than the cache holds, as concurrent requests
			// do. Reuse the distinct resources after draining each round's cache.
			resources := make([]interface{}, 32)
			for i := range resources {
				resources[i] = tc.acquire(b)
			}
			for round := 0; round < 5000; round++ {
				var wg sync.WaitGroup
				start := make(chan struct{})
				done := make(chan struct{})
				wg.Add(len(resources))
				for _, resource := range resources {
					resource := resource
					go func() {
						defer wg.Done()
						<-start
						tc.release(b, resource)
					}()
				}
				go func() { wg.Wait(); close(done) }()
				close(start)

				timer := time.NewTimer(time.Second)
				select {
				case <-done:
					timer.Stop()
				case <-timer.C:
					// Do not leak blocked senders on an unfixed implementation.
					cleanup := time.NewTimer(5 * time.Second)
					defer cleanup.Stop()
					for {
						select {
						case <-done:
							t.Fatalf("concurrent releases blocked with a full cache at round %d", round)
						case <-cleanup.C:
							t.Fatal("timed out cleaning up blocked releases")
						default:
							tc.drain(b)
							runtime.Gosched()
						}
					}
				}
				if !tc.drain(b) {
					t.Fatal("release did not retain any resource")
				}
				if tc.drain(b) {
					t.Fatal("retained more resources than the cache capacity")
				}
			}
		})
	}
}

func TestBoundedCachedCompressorsReleaseCapacity(t *testing.T) {
	for _, capacity := range []int{0, 1, 3} {
		for _, tc := range boundedCompressorCases() {
			t.Run(tc.name+"/"+strconv.Itoa(capacity), func(t *testing.T) {
				b := NewBoundedCachedCompressors(capacity, capacity)
				resources := make([]interface{}, capacity+2)
				for i := range resources {
					resources[i] = tc.acquire(b)
				}
				for _, resource := range resources {
					tc.release(b, resource)
				}
				for i := 0; i < capacity; i++ {
					if !tc.drain(b) {
						t.Fatalf("cache only retained %d resources, want %d", i, capacity)
					}
				}
				if tc.drain(b) {
					t.Fatal("release exceeded cache capacity")
				}
			})
		}
	}
}

func TestBoundedCachedCompressorsReusesReleasedResource(t *testing.T) {
	for _, tc := range boundedCompressorCases() {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBoundedCachedCompressors(1, 1)
			want := tc.acquire(b)
			for i := 0; i < 20; i++ {
				tc.release(b, want)
				if got := tc.acquire(b); got != want {
					t.Fatal("acquire did not reuse the released resource")
				}
			}
		})
	}
}
