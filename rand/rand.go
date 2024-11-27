package rand

import (
	"math/rand"
)

func GenRandNum(low, top int) int {
	return (rand.Intn(top-low) + low)
}

func GenRandString(length int) string {
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var result []rune
	for i := 0; i < length; i++ {
		j := GenRandNum(0, len(chars)-1)
		result = append(result, rune(chars[j]))
	}

	return string(result)
}
