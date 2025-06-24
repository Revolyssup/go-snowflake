package idgen

import (
	"sync"
	"testing"
)

func TestIDGenTimeSortable(t *testing.T) {
	worker := func(sf *Snowflake, wg *sync.WaitGroup, ch chan<- uint64) {
		defer wg.Done()
		ch <- sf.GetUUID()
	}
	sf := NewSnowflake(nil)
	var wg sync.WaitGroup
	ch := make(chan uint64, 1_000)
	for range 1_000_000 {
		wg.Add(1)
		go worker(sf, &wg, ch)
	}
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func(ch <-chan uint64) {
		defer wg2.Done()
		var curr uint64
		for val := range ch {
			if val <= curr {
				t.Error("Produced lower or equal value and not higher")
				t.Fail()
			}
		}
	}(ch)

	wg.Wait()
	close(ch)
	wg2.Wait()
}

func BenchmarkIDGen(b *testing.B) {
	sf := NewSnowflake(nil)
	for b.Loop() {
		sf.GetUUID()
	}
}
