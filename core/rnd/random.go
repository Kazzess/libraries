package rnd

import (
	"fmt"
	"math/rand"
	"time"
)

// RandInt returns integer random number
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())

	return min + rand.Intn(max-min)
}

// RandInt64 returns integer random number (int64).
func RandInt64(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())

	return min + rand.Int63n(max-min)
}

// RandFloat64 returns float random number.
func RandFloat64(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())

	return min + rand.Float64()*(max-min)
}

// RandomIP returns a random IP-address.
func RandomIP() string {
	return fmt.Sprintf(
		"%d.%d.%d.%d",
		RandInt(1, 255),
		RandInt(1, 255),
		RandInt(1, 255),
		RandInt(1, 255),
	)
}

// RandomDate returns a random date.
func RandomDate(min, max time.Time) time.Time {
	ts := RandInt64(min.Unix(), max.Unix())
	return time.Unix(ts, 0)
}

var (
	set    = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
	setLen = len(set) - 1
)

// RandomString returns a random string.
func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	randStr := make([]byte, n)

	for i := 0; i < n; i++ {
		r := rand.Intn(setLen)
		randStr[i] = set[r]
	}

	return string(randStr)
}

// RandomCase returns a random value from args.
func RandomCase(args ...interface{}) interface{} {
	return args[RandInt(0, len(args))]
}

// RandomBool returns a boolean random value.
func RandomBool() bool {
	return RandomCase(true, false).(bool)
}
