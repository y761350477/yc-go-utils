package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Int转String
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// String转Int
func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

// struct 转 map
func Struct2Map(obj interface{}) map[string]interface{} {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data[typeOfType.Field(i).Name] = field.Interface()
	}
	return data
}

// jsonStr转map
func JsonStrToMap(jsonstr string) map[string]string {
	var dat map[string]string
	if err := json.Unmarshal([]byte(jsonstr), &dat); err == nil {
		fmt.Println("==============json str 转map=======================")
		fmt.Println(dat)
		fmt.Println(dat["name"])
	}
	return dat
}

// struct转json
func StructToJsonOfByteArray(x interface{}) (bytes []byte, err error) {
	// 转换成JSON返回的是byte[]
	bytes, err = json.Marshal(x)
	if err != nil {
		fmt.Println(err.Error())
	}
	// byte[]转换成string 输出
	fmt.Println(string(bytes))
	return
}

// OK
func GbkToUtf8(s []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return "", e
	}
	return string(d), nil
}
