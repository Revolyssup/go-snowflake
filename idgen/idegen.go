package idgen

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

/*
Functional Requirement: Create unique time sortable 64 bit UUIDs for distributed systems.

Low level API: idgen.NewSnowflake(conf) -> Snowflake
snowflake.GetUUID() uint64
conf = {
instance_id 10 bit (no ranging from 0 to 1023)
custom_timestamp time.Time
}
42 bit timestamp
10 bit instance id
12 bit seq no

Either we can support more nodes or higher rate. I opt here for higher rate and give more bits to sequence number as a result.
Can generate 4096 unique IDs per millisecond per machine
snowflake.GetUUIDForTimestamp() (instead of time.Now converted to 41 bit unix ms, this will be converted)
Impl: using Twitter's Snowflake spec.
*/

type SnowflakeConfig struct {
	InstanceID      uint16
	CustomTimestamp *time.Time
}

type Snowflake struct {
	instanceId      uint16
	customtimestamp *time.Time
	sequenceNumber  uint32
	state           uint64
	mx              sync.Mutex
}

const (
	TimestampMask      = ((1 << 42) - 1)
	TimestampShift     = 64 - 42
	InstanceIDMask     = (1<<10 - 1)
	InstanceIDShift    = (64 - 42 - 10)
	SequenceNumberMask = (1<<12 - 1)
)

func (sf *Snowflake) setSeqNumber() {
	timer := time.NewTicker(time.Millisecond)
	for range timer.C {
		atomic.SwapUint32(&sf.sequenceNumber, 0)
	}
}
func NewSnowflake(conf *SnowflakeConfig) *Snowflake {
	if conf == nil {
		conf = &SnowflakeConfig{
			InstanceID: uint16(rand.Uint32()),
		}
	}

	sf := &Snowflake{
		instanceId:      conf.InstanceID,
		customtimestamp: conf.CustomTimestamp,
		mx:              sync.Mutex{},
	}
	go sf.setSeqNumber()
	return sf
}

func (sf *Snowflake) GetUUID() uint64 {
	var timestamp time.Time
	if sf.customtimestamp == nil {
		timestamp = time.Now()
	} else {
		timestamp = *sf.customtimestamp
	}

	return sf.GetUUIDForTimestamp(timestamp)
}

func (sf *Snowflake) GetUUIDForTimestamp(timestamp time.Time) uint64 {
BEGIN:
	oldstate := atomic.LoadUint64(&sf.state)
	var ts, id, sn uint64
	ts = (uint64(timestamp.UnixMilli()&TimestampMask) << TimestampShift)
	id = ((uint64(sf.instanceId)) & InstanceIDMask) << InstanceIDShift
	newSeq := atomic.AddUint32(&sf.sequenceNumber, 1)
	sn = (uint64(newSeq) & SequenceNumberMask)
	state := ts | id | sn
	if atomic.CompareAndSwapUint64(&sf.state, oldstate, state) {
		return state
	} else {
		time.Sleep(time.Millisecond) // To avoid spinlock
		goto BEGIN
	}
}
