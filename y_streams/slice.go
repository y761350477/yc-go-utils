package y_streams

import (
	"strings"
)

// TODO
// 如何去除换文件名中的所有特殊字符
func Replace(sli []string, s string) (bool, string) {
	for _, v := range sli {
		r := strings.ReplaceAll(s, v, ``)

		if strings.Contains(s, v) {
			return true, v
		}
	}
	return false, ""
}
