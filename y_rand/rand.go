package y_rand

import (
	"math/rand"
	"time"
)

// 小数区间的随机数
func RandFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// 整数区间的随机数
func RandInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandSleepMilli(min, max int) {
	randInt := RandInt(min, max)
	time.Sleep(time.Millisecond * time.Duration(randInt))
}
