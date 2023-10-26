package misc

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Str2list 将json风格的数组字符串s转换为字符串
func Str2list(s string) []string {
	return strings.Split(s[1:len(s)-1], ",")
}

// MinMaxInt 求一个int切片的最大最小值
func MinMaxInt(array []int) (min int, minIndex int, max int, maxIndex int) {
	max = array[0]
	min = array[0]
	for i, value := range array {
		if max < value {
			max = value
			maxIndex = i
		}
		if min > value {
			min = value
			minIndex = i
		}
	}
	return
}

// MinMaxUInt 求一个int切片的最大最小值
func MinMaxUInt(array []uint32) (min uint32, minIndex int, max uint32, maxIndex int) {
	max = array[0]
	min = array[0]
	for i, value := range array {
		if max < value {
			max = value
			maxIndex = i
		}
		if min > value {
			min = value
			minIndex = i
		}
	}
	return
}

type CanCompare interface {
	Less(CanCompare) bool
}

// MinMax 求一个CanCompare切片的最大最小值
func MinMax(array []CanCompare) (min CanCompare, minIndex int, max CanCompare, maxIndex int) {
	max = array[0]
	min = array[0]
	for i, value := range array {
		if max.Less(value) {
			max = value
			maxIndex = i
		}
		if min.Less(value) {
			min = value
			minIndex = i
		}
	}
	return
}

// GetTimeStr 返回当前的时间，格式为2006-01-02 15:04:05
func GetTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func InputStr(len int) string {
	re := ""
	for i := 0; i < len; i++ {
		re += "_"
	}
	for i := 0; i < len; i++ {
		re += "\b"
	}
	return re
}

func Input(hint string, len int, a ...interface{}) error {
	print(hint + InputStr(len))
	_, err := fmt.Scan(a...)
	return err
}

func Stop() {
	fmt.Printf("输入任意键继续...")
	ClearIOBuffer()
	b := make([]byte, 1)
	// 不知道为什么清空缓冲区后，还是有残留一个ascii为10的字符。。。但是goland里面调试时好的。。
	_, err := os.Stdin.Read(b)
	_, err = os.Stdin.Read(b)
	if err != nil {
		return
	}
}

func ClearIOBuffer() {
	myReader := bufio.NewReader(os.Stdin)
	myReader.Reset(os.Stdin)
	_, _ = myReader.ReadString('\n')
}

// IsLegalOutURL 判断是否是合法的外链
func IsLegalOutURL(url string) bool {
	if len(url) < 6 {
		return false
	}
	if strings.Index(url, "://") == -1 {
		return false
	}

	return true
}

// Clear 清空屏幕
func Clear() {
	cmd := exec.Command("cmd.exe", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		print("clear fail")
	}
}

// GetRealDir 返回当前可执行程序的真正目录，如果是软链接，则返回软链接指向的文件的目录
func GetRealDir() string {
	// 获取当前可执行程序的路径
	exePath, _ := os.Executable()
	// 获取当前可执行程序的目录
	exeDir := filepath.Dir(exePath)
	// 如果是软链接，获取软链接的路径
	if linkPath, err := os.Readlink(exePath); err == nil {
		// 如果是软链接，则获取软链接的绝对路径
		linkDir := filepath.Dir(linkPath)
		exeDir = linkDir
	}
	return exeDir
}

func InputString(notice string) string {
	fmt.Print(notice)
	var fileNames string
	_, err := fmt.Scanln(&fileNames)
	if err != nil {
		return ""
	}
	return fileNames
}

func InputStringWithSpace(notice string) []string {
	var msg string
	fmt.Print(notice)
	reader := bufio.NewReader(os.Stdin)
	msg, _ = reader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	msgList := strings.Split(msg, " ")
	return msgList
}

func InputInt(notice string) int {
	fmt.Print(notice)
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return 0
	}
	return choice
}

func ReadGBKFileLine(addr string) (string, error) {
	file, err := os.Open(addr)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := simplifiedchinese.GB18030.NewDecoder()
	returnStr := ""
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF && len(line) == 0 {
			break
		}
		// 去除末尾的\r\n
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")

		line, _, err = transform.String(decoder, line)
		if err != nil {
			return "", err
		}
		returnStr += line
	}
	return returnStr, nil
}

func WriteGBKFileLine(addr string, content string) error {
	file, err := os.Create(addr)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := simplifiedchinese.GB18030.NewEncoder()

	encodedLine, err := encoder.String(content)
	if err != nil {
		return err
	}
	_, err = file.WriteString(encodedLine + "\r\n")
	if err != nil {
		return err
	}

	return nil
}
