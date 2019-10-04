package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	time.Local = time.FixedZone("Beijing", 8*3600)
	// time.Local = time.UTC
	str := FormatTime(time.Now(), TimePrecision_Seconds)
	fmt.Println(str)
}
