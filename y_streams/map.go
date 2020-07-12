package y_streams

import (
	"sort"
)

// 获取最小值的KEY对应的VALUE值
func GetMinKeyValue(dataMap map[float64]string) string {
	var exeIp string
	result := map[float64]string{}
	var keys []float64
	// 压入各个数据
	result = dataMap
	// 得到各个key
	for key := range result {
		keys = append(keys, key)
	}
	// 给key排序，从小到大
	sort.Sort(sort.Float64Slice(keys))

	// 注意：遍历keys，而不是遍历map
	for _, key := range keys {
		exeIp = result[key]
		break
	}
	return exeIp
}
