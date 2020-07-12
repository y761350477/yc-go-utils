package y_regex

import (
	"regexp"
	"strings"

	"yc-go-utils/convert"
)

// 正则匹配
func RegexStr(text, regex string) string {
	re := regexp.MustCompile(regex)
	return re.FindString(text)
}

func RegexStrSub(text, regex string) [][]string {
	re := regexp.MustCompile(regex)
	// 只匹配第一个满足的条件
	return re.FindAllStringSubmatch(text, -1)
}

// 删除正则匹配的内容
func RegexRemoveFirst(src, regex string) string {
	re := regexp.MustCompile(regex)
	findString := re.FindString(src)
	return strings.Replace(src, findString, ``, 1)
}

// 删除正则匹配的内容（会把所有匹配的情况全删除）
func RegexRemove(src, regex string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllString(src, ``)
}

// func RegexKeep(src, regex, num string) string {
// 	regex := `.*?(\(.*?\)){` + num + `}`
// 	str_ = RegexRemove(src, regex)
// }

// 获取匹配到的组数
func RegexStrLen(text, regex string) int {
	re := regexp.MustCompile(regex)
	match := re.FindAllStringSubmatch(text, -1)
	return len(match)
}

// 获取正则匹配的状态
func RegexStatus(text, regex string) bool {
	re := regexp.MustCompile(regex)
	str := re.FindString(text)
	return str != ""
}

// 过滤无效名称
// 提取`名稱： (HD1080P H264)(Prestige)(118kbi00001)KANBi専属第1弾！...`的名称为`(Prestige)(118kbi00001)KANBi専属第1弾！...`
func RegexName(parentDirName string, RegexLabel string) (str_ string) {
	status := RegexStatus(parentDirName, RegexLabel)
	if status {
		size := RegexStrLen(parentDirName, RegexLabel)
		// 保留两个()
		i := size - 2
		toInt := convert.IntToString(i)
		regex := `.*?(\(.*?\)){` + toInt + `}`
		str_ = RegexRemove(parentDirName, regex)
	}
	return
}
