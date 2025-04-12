package redis

import (
	"github.com/go-redis/redis/v8"
)

const (
	Nil = redis.Nil
)

type (
	StringCmd          = redis.StringCmd
	StatusCmd          = redis.StatusCmd
	IntCmd             = redis.IntCmd
	BoolCmd            = redis.BoolCmd
	StringStringMapCmd = redis.StringStringMapCmd
	FloatCmd           = redis.FloatCmd
	StringSliceCmd     = redis.StringSliceCmd
	ZSliceCmd          = redis.ZSliceCmd
	ScanCmd            = redis.ScanCmd
)
