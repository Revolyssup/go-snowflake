package idgen

import (
	"testing"
	"time"
)

func TestIDGenTimeSortable(t *testing.T) {
	sf := NewSnowflake(nil)
	var curr uint64
	counter := 0
	for range 1_000_000 {
		new, err := sf.GetUUID()
		if err == ErrRateLimitExceeded {
			if counter != SequenceNumberMask {
				t.Errorf("counter should be %d but actual value %d", SequenceNumberMask, counter)
				t.FailNow()
				return
			}
			counter = 0
			time.Sleep(time.Millisecond)
			continue
		}
		counter++
		if new <= curr {
			t.Errorf("new value recieved %d for uuid is smaller or eq to previous %d", new, curr)
			t.FailNow()
		} else {
			curr = new
		}
	}
}

func TestAccuracy(t *testing.T) {
	sf := NewSnowflake(nil)
	tim := time.Now()
	uuid, _ := sf.GetUUIDForTimestamp(tim)
	tim2 := sf.GetTimestampForUUID(uuid)
	timms := tim.UnixMilli()
	tim2ms := tim2.UnixMilli()
	if timms != tim2ms {
		t.Errorf("Given time %d and recieved back %d", timms, tim2ms)
	}
}
func BenchmarkIDGen(b *testing.B) {
	sf := NewSnowflake(nil)
	for b.Loop() {
		sf.GetUUID()
	}
}
