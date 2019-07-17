package utils

import (
	"strconv"
	"strings"
)

// ParseIntArray parse a seperated string into a int64 array
func ParseIntArray(str string, sep string, base int, bitsize int) ([]int64, error) {
	strArray := strings.Split(str, sep)
	// transform string array to int array
	interArray := make([]int64, len(strArray))
	for i, s := range strArray {
		if integer, err := strconv.ParseInt(s, base, bitsize); err == nil {
			interArray[i] = integer
		} else {
			return nil, err
		}
	}
	return interArray, nil
}

func ParseInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
