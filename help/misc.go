package help

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/toophy/mahonia"
	"io"
	"math/rand"
	"os"
	"runtime"
	"strings"
)

// 判断文件/文件存在
func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 随机字符串
const rand_seed = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@#$%><-+"

func RandToken() []byte {
	str := make([]byte, 20)
	for i := 0; i < 20; i++ {
		str[i] = rand_seed[rand.Intn(len(rand_seed)-1)]
	}
	return str
}

// 获取上一层目录
func GetPreDir(dir string) string {
	dir3 := ""
	if runtime.GOOS == "windows" {
		dir2 := strings.LastIndexAny(dir, "\\")
		dir3 = dir[:dir2] + "\\"
	} else {
		dir2 := strings.LastIndexAny(dir, "/")
		dir3 = dir[:dir2] + "/"
	}
	return dir3
}

func Gbk2Utf8(src string) string {
	enc := mahonia.NewDecoder("gbk")
	return enc.ConvertString(src)
}

func Utf82Gbk(src string) string {
	enc := mahonia.NewEncoder("gbk")
	return enc.ConvertString(src)
}

//生成随机字符串
func RandStr(strlen int) string {

	data := make([]byte, strlen)
	var num int
	for i := 0; i < strlen; i++ {
		num = rand.Intn(57) + 65
		for {
			if num > 90 && num < 97 {
				num = rand.Intn(57) + 65
			} else {
				break
			}
		}
		data[i] = byte(num)
	}
	return string(data)
}

func GetMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func CopyFile(src, des string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return -1, err
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		return -1, err
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}
