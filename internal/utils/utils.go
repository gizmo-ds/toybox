package utils

import (
	"math/rand"
	"time"
)

func RandomRange(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}
