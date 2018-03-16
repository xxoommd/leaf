package util

import "strings"

// CheckFilename 检查资源名是否合法
//  格式: md5[32].ext
func CheckFilename(filename string) bool {
	s := strings.Split(filename, ".")
	if len(s) != 2 {
		return false
	}

	// MD5 length
	if len(s[0]) != 32 {
		return false
	}

	return true
}
