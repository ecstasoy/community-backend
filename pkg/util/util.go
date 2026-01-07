package util

import (
	"math/rand"
	"strconv"
	"time"
)

// RandomNumeric generates a random numeric string of the specified size.
// It is used to generate verification codes
func RandomNumeric(size int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if size < 0 {
		panic("size: " + strconv.Itoa(size) + " < 0")
	}
	value := ""
	for index := 0; index < size; index++ {
		value += strconv.Itoa(r.Intn(10))
	}

	return value
}

// EndOfDay returns a time.Time representing the end of the day (23:59:59) for the given time.
// It is used to set the expiration time for Redis keys
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}
