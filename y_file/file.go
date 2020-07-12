package y_file

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chenhg5/collection"
	"github.com/mholt/archiver"
	"github.com/prometheus/common/log"

	"yc-go-utils/y_regex"
)

// 读取文件内容(整体读取)
func ReadFile(filePath string) (string, error) {
	// 读取文件性能比较
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// 读取文件内容(逐行读取)
func ReadFileLine(filePath string) (string, error) {
	var str strings.Builder
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		result, err1 := readLine(r)
		if err1 != nil {
			if err1 == io.EOF {
				break
			} else {
				return str.String(), err1
			}
		}
		str.WriteString(result + "\n")
	}

	return str.String(), nil
}

// 读取文件内容(指定缓存大小)
func ReadFileBuffer(filePath string, bufferSize int) (string, error) {
	fileRead, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fileRead.Close()
	bufReader := bufio.NewReader(fileRead)
	buf := make([]byte, bufferSize)

	var str strings.Builder
	for {
		readNum, err := bufReader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == readNum {
			break
		}
		str.Write(buf[0:readNum])
	}
	return str.String(), nil
}

// 写入文件内容（追加）
func WriterFileAppend(outFile string, content string) (err error) {
	f, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	writer.WriteString(content)
	writer.Flush()
	return
}

// 写入文件内容（追加避免重复）
func WriterFileAppendNoRepeat(outFile string, content string) (err error) {
	source, _ := ReadFileBuffer(outFile, 2048)
	if !strings.Contains(source, content) {
		f, err1 := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err1 != nil {
			fmt.Printf("WriterFileAppendNoRepeat-err: %v\n", err)
			err = err1
			return
		}
		defer f.Close()
		writer := bufio.NewWriter(f)
		writer.WriteString(content)
		writer.Flush()
	}
	return
}

// 写入文件内容（清空）
func WriterFileTruncate(outFile string, content string) (err error) {
	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		return
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	writer.WriteString(content)
	writer.Flush()
	return
}

// 匹配包含指定后缀的行数据 - 逐行
func GetFileEqLine(filePath string, filterType string) (string, error) {
	filterType = strings.ToLower(filterType)
	var str strings.Builder
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		result, err1 := readLine(r)
		if err1 != nil {
			if err1 == io.EOF {
				break
			} else {
				return str.String(), err1
			}
		}
		if filterType != `` {
			ext := filepath.Ext(result)
			if strings.ToLower(ext) == filterType {
				str.WriteString(result + "\n")
			}
		} else {
			str.WriteString(result + "\n")
		}
	}

	return str.String(), nil
}

func ReadFileInfo01(filePath string) (r_ string, err error) {
	var RegexText = `\(.+?\)`
	var workName, rarType, videoType, magnetType string
	f, err0 := os.Open(filePath)
	if err0 != nil {
		err = err0
		return
	}
	defer func() {
		if err1 := f.Close(); err1 != nil {
			err = err1
			return
		}
	}()
	r := bufio.NewReader(f)
	for {
		result, err2 := readLine(r)
		// 判断读取是否为名称信息
		name := y_regex.RegexName(result, RegexText)
		if name != `` {
			workName = name + "\n"
		}

		suffix := filepath.Ext(result)
		switch suffix {
		case `.rar`:
			rarType += result + "\n"
		case `.mp4`, `.wmv`:
			videoType += result + "\n"
		default:
			if strings.Contains(result, `magnet:`) {
				magnetType += result + "\n"
			}
		}
		if err2 != nil {
			if err2 == io.EOF {
				break
			} else {
				err = err2
				return
			}
		}
	}
	switch {
	case videoType == `` && rarType == ``:
		r_ = workName + magnetType
	case videoType != ``:
		r_ = workName + videoType
	case videoType == ``:
		r_ = workName + rarType
	}
	return
}
func ReadFileInfo02(filePath string) (r_ string, err error) {
	var RegexText = `\(.+?\)`
	var info string
	f, err0 := os.Open(filePath)
	if err0 != nil {
		err = err0
		return
	}
	defer func() {
		if err1 := f.Close(); err1 != nil {
			err = err1
			return
		}
	}()
	r := bufio.NewReader(f)
	for {
		result, err2 := readLine(r)
		// 判断读取是否为名称信息
		name := y_regex.RegexName(result, RegexText)
		if name != `` {
			info = name + "\n"
		}

		switch suffix := filepath.Ext(result); suffix {
		case `.mp4`, `.wmv`, `.rar`:
			info += result + "\n"
		default:
			if strings.Contains(result, `magnet:`) {
				info += result + "\n"
			}
		}
		if err2 != nil {
			if err2 == io.EOF {
				break
			} else {
				err = err2
				return
			}
		}
	}
	r_ = info
	return
}

// 读取文件,判断是否包含指定内容(逐行、bufio)
// 1. 忽略大小写
// 2. content为空算包含
func Contain(filePath string, content string) (bool, error) {
	if content == `` {
		return true, nil
	}
	if content != `` {
		f, err0 := os.Open(filePath)
		if err0 != nil {
			return false, err0
		}
		defer f.Close()
		r := bufio.NewReader(f)
		for {
			result, err2 := readLine(r)
			// 转小写比较
			contains := strings.Contains(strings.ToLower(result), strings.ToLower(content))
			if contains {
				return true, nil
			}
			if err2 != nil {
				if err2 == io.EOF {
					break
				} else {
					return false, err2
				}
			}
		}
	}
	return false, nil
}

func readLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}

// 下载需要访问外网、用户认证的文件
func DownloadFileCookie(fileUrl string, fileName string, cookie string, userAgent string) (err error) {

	req, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		fmt.Println("获取地址错误")
		return err
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
		req.Header.Set("User-Agent", userAgent)
	}
	// 代理
	// proxyUrl := "http://127.0.0.1:10808"
	// proxy, _ := url.Parse(proxyUrl)
	// tr := &http.Transport{
	// 	Proxy:           http.ProxyURL(proxy),
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	client := &http.Client{
		// Transport: tr,
		Timeout: time.Second * 6, // 超时时间
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("下载文件异常: %v\n", err)
		return err
	}

	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("file 创建文件异常: %v\n", err)
		return err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("file 下载文件异常: %v\n", err)
		return err
	}
	return
}

func DownloadFile(fileUrl string, fileName string, proxyUrl string) (err error) {
	req, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		fmt.Println("获取地址错误")
		return
	}

	// 代理
	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 90, // 超时时间
	}

	resp, err := client.Do(req)

	if err != nil {
		// fmt.Printf("下载文件异常: %v\n", err)
		return
	}

	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create(fileName)
	defer out.Close()

	if err != nil {
		fmt.Printf("file 创建文件异常: %v\n", err)
		return
	}
	pix, err := ioutil.ReadAll(resp.Body)
	_, err = io.Copy(out, bytes.NewReader(pix))
	if err != nil {
		fmt.Printf("file 读取文件异常: %v\n", err)
		return
	}
	if err != nil {
		fmt.Printf("file 下载文件异常: %v\n", err)
		return
	}
	return
}

// 判断文件夹是否存在, 不存在则创建
func MkdirFolder(path string) (err error) {
	exists, err := PathExists(path)
	if !exists {
		os.MkdirAll(path, os.ModePerm)
	}
	return
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 判断文件是否存在,不存在则创建
func MkdirFile(filePath string) (err error) {
	dir := filepath.Dir(filePath)
	err = MkdirFolder(dir)
	f, err := os.OpenFile(filePath, os.O_CREATE, 0644)
	defer f.Close()
	return
}

// 获取指定目录下的所有文件和目录
func GetFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetFilesAndDirs(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".txt")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	return files, dirs, nil
}

// 获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string, filter string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth+PthSep+fi.Name(), filter)
		} else {
			// 过滤指定格式
			if filter != "" {
				ok := strings.HasSuffix(fi.Name(), filter)
				if ok {
					files = append(files, dirPth+PthSep+fi.Name())
				}
			} else {
				files = append(files, dirPth+PthSep+fi.Name())
			}

		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table, filter)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}

func GetAllDirs(dirPath string) (dirs []string, err error) {
	var tempDirs []string
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			dirPath1 := dirPath + PthSep + fi.Name()
			tempDirs = append(tempDirs, dirPath1)
			dirs = append(dirs, dirPath1)
			GetAllDirs(dirPath1)
		}
	}

	// 读取子目录下文件
	for _, d := range tempDirs {
		temp, _ := GetAllDirs(d)
		for _, temp1 := range temp {
			dirs = append(dirs, temp1)
		}
	}
	return
}

// 删除目录下所有空文件夹(递归遍历内部)
func DeleteTempDirs(dirPth string) (err error) {
	allDirs, err := GetAllDirs(dirPth)
	if err != nil {
		fmt.Printf("file.go : %v\n", err)
	}
	sort.StringSlice(allDirs).Sort()
	for i := len(allDirs) - 1; i >= 0; i-- {
		dir, err := ioutil.ReadDir(allDirs[i])
		if err != nil {
			return err
		}

		dirSize := len(dir)
		if dirSize <= 0 {
			os.Remove(allDirs[i])
		}
	}

	return
}

func GetPathDir(text string) string {
	return filepath.Dir(text)
}

func GetPathBase(text string) string {
	return filepath.Base(text)
}

// 判断RAR文件是否损坏
func ReadRarFileStatus(filePath string) (t_ bool) {
	r := archiver.NewRar()
	err := r.Walk(filePath, func(f archiver.File) error {
		return nil
	})
	if err != nil {
		t_ = true
	}
	return
}

func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func DeleteFile(path string) (err error) {
	if err := os.Remove(path); err != nil {
		log.Error(err)
	}
	return
}

// localPath指定删除目录
// outermost指定最外层目录
func DeleteFileOnDisk(localPath string, outermost string) {
	log.Debugf("remove file: %s", localPath)
	if err := os.Remove(localPath); err != nil {
		log.Error(err)
	}
	dirsList := make([]string, 0, 0)
	for dir := filepath.Dir(localPath); dir != outermost && len(dir) > len(outermost); dir = filepath.Dir(dir) {
		dirsList = append(dirsList, dir)
	}
	sort.StringSlice(dirsList).Sort()
	for i := len(dirsList) - 1; i >= 0; i-- {
		f, err := os.Open(dirsList[i])
		if err != nil {
			log.Error(err)
		}
		fs, err2 := f.Readdirnames(1)
		if err2 == io.EOF && (fs == nil || len(fs) == 0) {
			f.Close()
			log.Debugf("remove dir: %s", dirsList[i])
			if err := os.Remove(dirsList[i]); err != nil {
				log.Error(err)
			}
			continue
		} else if err2 != nil {
			log.Error(err2)
		}
		f.Close()
	}
}

// 判断目录下是否包含指定文件
func IsExist(path string, str string) (b bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		log.Error(err)
		return b, err
	}
	fs, err := f.Readdirnames(1)
	if err != nil {
		log.Error(err)
		return b, err
	}
	b = collection.Collect(fs).Contains(str)
	return
}

// 获取上级目录文件名
func GetParentDirName(path string) string {
	dir := filepath.Dir(path)
	return y_regex.RegexRemove(dir, `.*\\+`)
}

// 获取去除文件后缀的目录路径
func GetPathRemoveExt(text string) string {
	fileSuffix := path.Ext(text)
	return strings.TrimSuffix(text, fileSuffix)
}

// 是否是指定后缀的文件
func IsSuffixExt(text, suffix string) bool {
	fileSuffix := path.Ext(text)
	return suffix == fileSuffix
}
