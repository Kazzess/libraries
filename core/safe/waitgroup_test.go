package safe

import (
	"sync"
	"testing"
)

func TestWaitGroup_Add(t *testing.T) {
	wg := NewWaitGroup(nil)
	n := 100

	w := sync.WaitGroup{}
	w.Add(2)

	go func() {
		defer w.Done()

		for i := 0; i < n; i++ {
			go func(k int) {
				wg.Add(1)

				defer wg.Done()

				_ = k*k - k
			}(i)
		}
	}()

	go func() {
		defer w.Done()
		wg.Wait()
	}()

	w.Wait()
}

func BenchmarkWaitGroup_Add(b *testing.B) {
	wg := NewWaitGroup(nil)

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		wg.Done()
	}
}

func BenchmarkWaitGroup_AddParallel(b *testing.B) {
	wg := NewWaitGroup(nil)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg.Add(1)
			wg.Done()
		}
	})
}
