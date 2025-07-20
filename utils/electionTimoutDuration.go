package utils

import (
	"math/rand"
	"time"
)

func ElectionTimeoutDuration() time.Duration {

	randomMilli := rand.Intn(151) + 150
	duration := time.Duration(randomMilli) * time.Millisecond

	return duration
}
