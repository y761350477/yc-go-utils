package y_time

import (
	"fmt"
	"time"

	"yc-go-utils/y_file"
)

// 获取日期时间
func NowDateTime() (str_ string) {
	str_ = time.Now().Format("2006-01-02 15:04:05")
	return str_
}

// 获取当前日期
func NowDate() (str_ string) {
	str_ = time.Now().Format("2006-01-02")
	return str_
}

// 获取当前时间
func NowTime() (str_ string) {
	str_ = time.Now().Format("15:04:05")
	return str_
}

// 求时间差
func TimeDiff(f func()) time.Duration {
	startTime := time.Now()
	f()
	endTime := time.Now()
	return endTime.Sub(startTime)
}

// 指定睡眠时间
func SleepSecond(i int) {
	time.Sleep(time.Second * time.Duration(i))
}

// 写入操作时间
func WriterFileTime(path string) {
	bool, _ := y_file.Contain(path, NowDate())
	timeStyle := TimeStyle(bool)
	if err := y_file.WriterFileAppend(path, timeStyle); err != nil {
		fmt.Printf("writerFileTime-err: %v\n", err)
		return
	}
}

// 先判断今天是否操作过?
// 操作过，添加当前时间
// 没操作过，添加当前日期
func TimeStyle(b bool) (s string) {
	if b {
		s = "------------------------------------      Time: " + NowTime() + "        ------------------------------------\n"
	} else {
		s = "====================================      Date: " + NowDate() + "      ====================================\n" +
			"------------------------------------      Time: " + NowTime() + "        ------------------------------------\n"
	}
	return
}
