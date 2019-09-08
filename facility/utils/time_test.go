package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	str := FormatTime(time.Now(), TimePrecision_Seconds)
	fmt.Println(str)
}
