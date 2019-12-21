package utils

import "math/rand"

func GenerateId(length int) int64 {
	var res int64 = 0
	res = int64(rand.Intn(9)) + 1
	for i := 0; i < length-1; i++ {
		res = res*10 + int64(rand.Intn(10))
	}
	return res
}
