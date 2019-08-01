package strings

const EMPTY = ""

func PtrString(str string) *string {
	return &str
}

func PtrValue(ptr *string) string {
	if ptr == nil {
		return EMPTY
	}
	return *ptr
}

func IsEmpty(str string) bool {
	return str == EMPTY
}
