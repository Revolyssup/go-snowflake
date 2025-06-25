package idgen

import (
	"testing"
)

func TestIDGenTimeSortable(t *testing.T) {
	sf := NewSnowflake(nil)
	var curr uint64
	for range 1_000_0 {
		new := sf.GetUUID()
		if new <= curr {
			t.Errorf("new value recieved %d for uuid is smaller or eq to previous %d", new, curr)
			t.FailNow()
		} else {
			curr = new
		}
	}
}

func BenchmarkIDGen(b *testing.B) {
	sf := NewSnowflake(nil)
	for b.Loop() {
		sf.GetUUID()
	}
}
