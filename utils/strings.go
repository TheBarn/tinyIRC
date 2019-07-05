package utils

//StringInSlice will return true if the slice contains the string
func StringInSlice(str string, sl []string) bool {
	for _, val := range sl {
		if str == val {
			return true
		}
	}
	return false
}
