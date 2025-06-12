package Tools

import (
	"bytes"
	"compress/gzip"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

//var httpClient = &http.Client{
//	Timeout: time.Second * 2,
//}

// 全局 transport
var transport = &http.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 20,
	IdleConnTimeout:     30 * time.Second,
	MaxConnsPerHost:     100,
	DisableKeepAlives:   false,
	//TLSHandshakeTimeout:   5 * time.Second,
	//ResponseHeaderTimeout: 10 * time.Second,
}

var httpClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: transport,
}

//func init() {
//	go func() {
//		for {
//			fmt.Println("定时清理http链接")
//			time.Sleep(time.Minute * 1)
//			transport.CloseIdleConnections()
//		}
//	}()
//}

// POST josn请求工具 返回body 结构
func PostJsonRequestTool(url string, jsonMap map[string]interface{}, headerList map[string]string) (map[string]interface{}, error) {
	//// 创建自定义 Transport 用于连接池管理
	//transport := &http.Transport{
	//	MaxIdleConns:        10,               // 最大空闲连接数
	//	IdleConnTimeout:     30 * time.Second, // 空闲连接的超时时间
	//	DisableKeepAlives:   false,            // 启用连接池，允许 KeepAlive
	//	MaxIdleConnsPerHost: 10,               // 每个主机的最大空闲连接数
	//}
	//
	//// 创建 HttpClient，并使用自定义 Transport
	//httpClient := &http.Client{
	//	Transport: transport,
	//	Timeout:   10 * time.Second, // 请求超时时间
	//}

	jsonByte, err1 := jsoniter.Marshal(jsonMap) // 转换请求参数map为byte[]
	if err1 != nil {
		return nil, err1
	}

	var jsonString string
	jsonString = strings.ReplaceAll(string(jsonByte), `\u0026`, "&")
	//jsonString = strings.ReplaceAll(jsonString, `/`, "\\/")
	bf := bytes.NewBufferString(jsonString)

	request, err2 := http.NewRequest("POST", url, bf)
	if err2 != nil {
		return nil, err2
	}

	request.Close = true
	request.Header.Add("Connection", "close")

	if headerList == nil {
		request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	} else {
		for hkey, hvalue := range headerList {
			request.Header.Add(hkey, hvalue)
		}
	}

	response, err3 := httpClient.Do(request)
	if err3 != nil {
		if response != nil {
			response.Body.Close()
		}
		return nil, err3
	}

	//defer response.Body.Close()
	//defer transport.CloseIdleConnections()

	defer func() {
		err := response.Body.Close()
		if err != nil {
			for i := 0; i < 3; i++ {
				err = response.Body.Close()
				if err == nil {
					break
				}
			}
			if err != nil {
				log.Error(fmt.Sprintf("response关闭失败，%v", err))
			}
		}
	}()

	var body []byte
	var err4 error
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err4_1 := gzip.NewReader(response.Body)
		if err4_1 != nil {
			if reader != nil {
				reader.Close()
			}
			return nil, err4_1
		}
		//defer reader.Close()

		defer func() {
			err := reader.Close()
			if err != nil {
				for i := 0; i < 3; i++ {
					err = reader.Close()
					if err == nil {
						break
					}
				}
				if err != nil {
					log.Error(fmt.Sprintf("reader关闭失败，%v", err))
				}
			}
		}()

		body, err4 = io.ReadAll(reader)
		if err4 != nil {
			return nil, err4
		}
	} else {
		body, err4 = io.ReadAll(response.Body)
		if err4 != nil {
			return nil, err4
		}
	}

	bodyStr := string(body)
	bodyDst := make(map[string]interface{})
	err5 := jsoniter.Unmarshal([]byte(bodyStr), &bodyDst)
	if err5 != nil {
		return nil, err5
	}

	return bodyDst, nil
}

// POST	formdata
func PostFormDataRequestTool(requrl string, jsonStr string, headerList map[string]string) (map[string]interface{}, error) {

	//jsonData := url.Values{}
	//
	////jsonData := make(url.Values)
	//for jkey, jvalue := range jsonMap {
	//	fmt.Println(jkey, jvalue)
	//	jsonData.Add(jkey, jvalue.(string))
	//}
	//
	//fmt.Println(jsonData)

	//request, err2 := http.NewRequest("POST", requrl, strings.NewReader(jsonData.Encode()))
	request, err2 := http.NewRequest("POST", requrl, strings.NewReader(jsonStr))

	if err2 != nil {
		return nil, err2
	}

	request.Close = true
	request.Header.Set("Connection", "close")

	if headerList == nil {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		for hkey, hvalue := range headerList {
			request.Header.Add(hkey, hvalue)
		}
	}

	response, err3 := httpClient.Do(request)

	if err3 != nil {
		if response != nil {
			response.Body.Close()
		}
		return nil, err3
	}

	defer response.Body.Close()
	//defer transport.CloseIdleConnections()

	var body []byte
	var err4 error
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err4_1 := gzip.NewReader(response.Body)
		if err4_1 != nil {
			if reader != nil {
				reader.Close()
			}
			return nil, err4_1
		}
		defer reader.Close()

		body, err4 = io.ReadAll(reader)
		if err4 != nil {
			return nil, err4
		}
	} else {
		body, err4 = io.ReadAll(response.Body)
		if err4 != nil {
			return nil, err4
		}
	}

	bodyStr := string(body)
	bodyDst := make(map[string]interface{})
	err5 := jsoniter.Unmarshal([]byte(bodyStr), &bodyDst)
	if err5 != nil {
		return nil, err5
	}

	return bodyDst, nil
}

// Get请求工具 返回body 结构
func GetRequestTool(url string, headerList map[string]string) (map[string]interface{}, error) {
	//payload := &bytes.Buffer{}
	//writer := multipart.NewWriter(payload)
	//err := writer.Close()
	//if err != nil {
	//	return nil, err
	//}

	request, err2 := http.NewRequest("GET", url, nil)
	if err2 != nil {
		return nil, err2
	}

	request.Close = true
	request.Header.Set("Connection", "close")

	for hkey, hvalue := range headerList {
		request.Header.Add(hkey, hvalue)
	}

	//client := &http.Client{
	//	Timeout: time.Second * 10,
	//}

	response, err3 := httpClient.Do(request)

	if err3 != nil {
		if response != nil {
			response.Body.Close()
		}
		return nil, err3
	}

	defer response.Body.Close()
	//defer transport.CloseIdleConnections()

	var body []byte
	var err4 error
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err4_1 := gzip.NewReader(response.Body)
		if err4_1 != nil {
			if reader != nil {
				reader.Close()
			}
			return nil, err4_1
		}
		defer reader.Close()

		body, err4 = io.ReadAll(reader)
		if err4 != nil {
			return nil, err4
		}
	} else {
		body, err4 = io.ReadAll(response.Body)
		if err4 != nil {
			return nil, err4
		}
	}

	bodyStr := string(body)
	bodyDst := make(map[string]interface{})
	err5 := jsoniter.Unmarshal([]byte(bodyStr), &bodyDst)
	if err5 != nil {
		return nil, err5
	}

	return bodyDst, nil
}

// 获取长链接
func GetLongUrl(url string) (string, error) {

	//payload := &bytes.Buffer{}
	//writer := multipart.NewWriter(payload)
	//defer writer.Close()
	//err := writer.Close()
	//if err != nil {
	//	return "", err
	//}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	req.Close = true
	req.Header.Set("Connection", "close")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	//defer transport.CloseIdleConnections()

	return res.Request.URL.String(), nil
}
