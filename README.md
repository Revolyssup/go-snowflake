# go-snowflake

Implementation of twitter's snowflake algorithm for generating UUIDs for distributed systems.
Inspired from: https://github.com/tangledbytes/nodejs-snowflake/blob/master/README.md


## Usage

```go
snowflake := idgen.NewSlowfake(nil)
id := snowflake.GetUUID()
```

For generating ID with custom instance id and/or timestamp

```go
snowflake := idgen.NewSlowfake(&idgen.SnowflakeConfig{
	InstanceID: 23 //can be a number betwen 0 and 1023. For bigger numbers, the higher bit values are ignored
	CustomTimestamp: time.Now()
})
id := snowflake.GetUUID()
```

```note
The implementation is go routine safe
```
