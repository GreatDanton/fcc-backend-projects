package utilities

// StringInSlice returns true if string exists in array
// and false if it does not
func StringInSlice(str string, arr []string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
