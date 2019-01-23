package test

import (
	"strconv"
	"strings"
	"testing"
)

func TestStringSplit(t *testing.T) {
	var str = "1,2,3,4"
	splittedStr := strings.Split(str, ",")
	t.Logf("splitted string is %v", splittedStr)
	// transform string array to int array
	interArray := make([]int64, len(splittedStr))
	for i, s := range splittedStr {
		if integer, err := strconv.ParseInt(s, 10, 64); err == nil {
			interArray[i] = integer
		}
	}
	t.Logf("integer array is %v", interArray)
}
