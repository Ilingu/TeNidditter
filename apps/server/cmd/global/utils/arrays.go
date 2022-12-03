package utils

import (
	"math/rand"
	"time"
)

// Takes a pointer to the array and suffle it randomly
func ShuffleSlice[T any](slice []T) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
