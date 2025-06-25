# go-snowflake

Implementation of twitter's snowflake algorithm for generating UUIDs for distributed systems.


## Usage

```go
snowflake := idgen.NewSlowfake(nil)
id, err := snowflake.GetUUID()
```

For generating ID with custom instance id and/or timestamp

```go
snowflake := idgen.NewSlowfake(&idgen.SnowflakeConfig{
	InstanceID: 23 //can be a number betwen 0 and 1023. For bigger numbers, the higher bit values are ignored
	CustomTimestamp: time.Now()
})
id, err := snowflake.GetUUID()
```

```note
- The implementation is go routine safe
- The error returned is "Rate limit exceeded" and it is upto the client to handle this case.
```

## Benchmark

```bash
 go test -bench=. ./idgen
goos: linux
goarch: amd64
pkg: github.com/Revolyssup/go-snowflake/idgen
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkIDGen-16       22185558                45.39 ns/op
PASS
```

### Reasoning behind throwing error in case of rate limit as opposed to internal retries

- If for a given timestamp the sequence number has overflown then the only way to reliably return an id is by incrementing the timestamp+1 but I think
this is non intuitive and inaccurate.
- Blocking requests in case of rate limit hides the behavious from client and in case of uncontrolled requests might lead to unexpected behavior.
- Allowing clients to control what happens in case of rate limit allows more control to clients: They may have their own retry mechanism which is best suited
to their system.
