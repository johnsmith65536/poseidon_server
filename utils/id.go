package utils

import "math/rand"

const characterSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func GenerateId(length int) int64 {
	var res int64 = 0
	res = int64(rand.Intn(9)) + 1
	for i := 0; i < length-1; i++ {
		res = res*10 + int64(rand.Intn(10))
	}
	return res
}

func GenerateUUID() string {
	var res string
	for i := 1; i <= 16; i++ {
		var x = rand.Intn(62)
		res += characterSet[x : x+1]
		if i%4 == 0 && i != 16 {
			res += "-"
		}
	}
	return res
}

func GenerateToken(length int) string {
	var token string
	for i := 0; i < length; i++ {
		var x = rand.Intn(62)
		token += characterSet[x : x+1]
	}
	return token
}
