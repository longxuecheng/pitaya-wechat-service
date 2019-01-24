package utils

// CheckAndPanic check whether the err is null and throw it
func CheckAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}
