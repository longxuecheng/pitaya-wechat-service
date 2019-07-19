package strings

const EMPTY = ""

func PtrString(str string) *string {
	return &str
}

func IsEmpty(str string) bool {
	return str == EMPTY
}
