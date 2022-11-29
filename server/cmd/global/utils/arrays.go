package utils

import (
	"math/rand"
	"time"
)

func ShuffleSlice[T any](slice []T) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
