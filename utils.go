package utils

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"text/template"

	"yc-go-utils/convert"
)

// 变量替换
func VariableReplace(enContent string, dynamicVar string) string {
	tpl, err := template.New("tplReplace").Parse(enContent)
	if err != nil {
		panic(err.Error())
	}

	tmpMap := convert.JsonStrToMap(dynamicVar)
	var buf bytes.Buffer
	err = tpl.Execute(&buf, tmpMap)
	if err != nil {
		panic(err.Error())
	}

	return buf.String()
}

// 获取ip
func GetThisIp() string {
	var ipStr string
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ipStr = ipnet.IP.String()
					}
				}
			}
		}
	}
	return ipStr
}

func GetAppList(name string) (result *bufio.Scanner, err error) {
	fptr := flag.String("fpath", name, "启动列表路径")
	flag.Parse()
	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	s := bufio.NewScanner(f)
	err = s.Err()
	return s, err
}
