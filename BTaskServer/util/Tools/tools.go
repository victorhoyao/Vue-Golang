package Tools

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"reflect"
	"strconv"
	"time"
)

// md5加密
func GenMd5(code string) string {
	Md5 := md5.New()
	io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

// 生成uuid
func GetUuid() string {
	id := uuid.NewV4()
	return id.String()
}

// 生成随机数
func GetRand(max string) string {
	num, err := strconv.Atoi(max)
	if err != nil {
		return max
	}
	max_len := strconv.Itoa(len(max))
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(num)
	n1 := fmt.Sprintf("%0"+max_len+"d", n)
	return n1
}

// 文件保存到服务器本地
func SaveFileByLocal(file *multipart.FileHeader, savepath, savename string) bool {
	f1, err := os.OpenFile(savepath+"//"+savename, os.O_CREATE|os.O_WRONLY, 0777) // 不存在创建,只写打开
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer f1.Close()

	f, err := file.Open()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(f)
	b := make([]byte, 1024)
	for {
		n, err := reader.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return false
		}
		f1.Write(b[0:n])
	}
	return true
}

// 获取当前的日期 格式化
func GetDateNowFormat(flag bool) string {
	if flag {
		template := "2006-01-02 15:04:05"
		return time.Now().Format(template)
	} else {
		template := "2006-01-02"
		return time.Now().Format(template)
	}
}

// 获取秒级时间戳
func GetNowUnix(b bool) int64 {
	if b == true {
		return time.Now().Unix()
	} else {
		return time.Now().UnixNano()
	}
}

func Struct2map(obj any) (data map[string]any, err error) {
	// 通过反射将结构体转换成map
	data = make(map[string]any)
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	for i := 0; i < objT.NumField(); i++ {
		fileName, ok := objT.Field(i).Tag.Lookup("bson")
		if ok {
			data[fileName] = objV.Field(i).Interface()
		} else {
			data[objT.Field(i).Name] = objV.Field(i).Interface()
		}
	}
	return data, nil
}

func IsLaterNow(timestr string) bool {
	timeNow, _ := time.Parse("2006-01-02 15:04:05", GetDateNowFormat(true))
	timea, _ := time.Parse("2006-01-02 15:04:05", timestr)
	return timea.Before(timeNow)
}
