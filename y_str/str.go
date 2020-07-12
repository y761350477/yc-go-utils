package y_str

import (
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/chenhg5/collection"
)

const (
	// 路径分隔符
	PthSep = string(os.PathSeparator)

	RegexText = `\(.+?\)`
)

// 指定字符分隔
func StrSplit(src, str string) []string {
	result := strings.Split(src, str)
	sle := make([]string, 0, len(result))
	for _, v := range result {
		v = strings.TrimSpace(v)
		sle = append(sle, v)
	}
	return sle
}

// 获取当前条数在第几页（根据每页总条数，当前条数计算）
// num当前条数
// pageTotal每页总条数
func GetPageNum(num, pageTotal int) int {
	i_ := num / pageTotal
	i := num % pageTotal
	if i == 0 {
		return i_
	}
	return i_ + 1
}

func RemoveTexts(src string, texts []string) string {
	for _, v := range texts {
		src = strings.ReplaceAll(src, v, ``)
	}
	return src
}

// 获取优化后作品的名称
// func GetLabelName(name string) string {
// 	return RegexName(name, `\(.+?\)`)
// }

// 判断slice中是否存在某个item
func IsExistItem(src interface{}, item interface{}) bool {
	return collection.Collect(src).Contains(item)
}

// 复制内容到剪切板
func ClipText(text string) {
	clipboard.WriteAll(text)
}

// 获取文件内容首行信息
// func GetContentFirst(text string) string {
// 	return RegexStr(text, `.*\n`)
// }

// 提取番号
// func GetSd(src string, RegexText1 string, spli []string) (str_ string) {
// 	re := regexp.MustCompile(src)
// 	f := re.FindAllString(RegexText1, -1)
// 	if len(f) == 1 {
// 		str_ = f[0]
// 		return
// 	} else {
// 		for _, s_v := range spli {
// 			for _, f_v := range f {
// 				if s_v != f_v {
// 					str_ = f_v
// 					return
// 				}
// 			}
// 		}
// 	}
// 	return
// }

// 获取正则匹配的状态
func RegexStatus(text, regex string) bool {
	re := regexp.MustCompile(regex)
	str := re.FindString(text)
	return str != ""
}

// 根据分隔符分割字符串生成切片
// func StrSplit(data, str string) []string {
// 	result := strings.Split(data, str)
// 	sle := make([]string, 0, len(result))
// 	for _, v := range result {
// 		v = strings.TrimSpace(v)
// 		sle = append(sle, v)
// 	}
// 	return sle
// }
