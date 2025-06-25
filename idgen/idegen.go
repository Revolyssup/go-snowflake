package idgen

import (
	"fmt"
	"math/rand"
	"sync"
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
Can generate 4096 unique IDs per millisecond per machine.
For simplicity, accuracy, predictability and speed the function fast-fails when rate limit is exceeded allowing clients to decide the next move.
*REASONING BEHIND IT:
*If for a given timestamp the sequence number has overflown then the only way to reliably return an id is by incrementing the timestamp+1 but I think
*this is non intuitive and inaccurate.

*Blocking requests in case of rate limit hides the behavious from client and in case of uncontrolled requests might lead to unexpected behavior.
*Allowing clients to control what happens in case of rate limit allows more control to clients: They may have their own retry mechanism which is best suited
*to their system.
*/
var ErrRateLimitExceeded = fmt.Errorf("rate limit exceeded. try after a millisecond")

type SnowflakeConfig struct {
	InstanceID      uint16
	CustomTimestamp *time.Time
}

type Snowflake struct {
	instanceId      uint16
	customtimestamp *time.Time
	sequenceNumber  uint32
	lasttimestamp   uint64
	mx              sync.Mutex
}

const (
	TimestampMask      = ((1 << 42) - 1)
	TimestampShift     = 64 - 42
	InstanceIDMask     = (1<<10 - 1)
	InstanceIDShift    = (64 - 42 - 10)
	SequenceNumberMask = (1<<12 - 1)
)

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
	return sf
}

func (sf *Snowflake) GetUUID() (uint64, error) {
	var timestamp time.Time
	if sf.customtimestamp == nil {
		timestamp = time.Now()
	} else {
		timestamp = *sf.customtimestamp
	}

	return sf.GetUUIDForTimestamp(timestamp)
}

func (sf *Snowflake) GetUUIDForTimestamp(timestamp time.Time) (uint64, error) {
	sf.mx.Lock()
	defer sf.mx.Unlock()
	var ts, id, sn uint64
	ts = (uint64(timestamp.UnixMilli()&TimestampMask) << TimestampShift)
	defer func() {
		sf.lasttimestamp = ts
	}()
	id = ((uint64(sf.instanceId)) & InstanceIDMask) << InstanceIDShift
	sn = (uint64(sf.sequenceNumber) & SequenceNumberMask)
	if sn == SequenceNumberMask { //overflow
		sf.sequenceNumber = 0
		if ts == sf.lasttimestamp {
			return 0, ErrRateLimitExceeded
		}
	} else {
		sf.sequenceNumber++
	}
	state := ts | id | sn
	return state, nil
}

// This will return an approximate value of time with ms accuracy
func (sf *Snowflake) GetTimestampForUUID(uuid uint64) time.Time {
	ts := (uuid >> TimestampShift)
	return time.UnixMilli(int64(ts))
}
