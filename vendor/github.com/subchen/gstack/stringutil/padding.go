package stringutil

import "strings"

func LeftPad(s string, padStr string, size int) string {
	var padCountInt int
	padCountInt = 1 + ((size - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - size):]
}

func RightPad(s string, padStr string, size int) string {
	var padCountInt int
	padCountInt = 1 + ((size - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:size]
}
