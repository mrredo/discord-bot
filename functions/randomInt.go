package functions

import (
	"math/rand"
	"time"
)

func RandomInRange(a, b int) int {
	// Initialize the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number in the specified range
	return rand.Intn(b-a+1) + a
}
