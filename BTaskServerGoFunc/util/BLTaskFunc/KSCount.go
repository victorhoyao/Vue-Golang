package BLTaskFunc

import (
	"BTaskServer/util/Tools"
	"fmt"
	"strconv"
)

func GetKSCount(shortUrl string, taskType int) map[string]interface{} {
	resMap := make(map[string]interface{})

	var url string
	if taskType == 11 {
		url = fmt.Sprintf("http://111.180.204.175:9999/api/index/getFans?url=%s", shortUrl)
	} else if taskType == 12 {
		url = fmt.Sprintf("http://111.180.204.175:9999/api/index/getLike?url=%s", shortUrl)
	}

	bodyDst, err := Tools.GetRequestTool(url, nil)

	if err != nil {
		resMap["ret"] = false
		resMap["msg"] = "请求错误"
		resMap["count"] = 0
		return resMap
	}

	if bodyDst["code"] == nil {
		resMap["ret"] = false
		resMap["msg"] = "code不存在"
		resMap["count"] = 0
		return resMap
	}

	code := int(bodyDst["code"].(float64))
	if code != 1 {
		resMap["ret"] = false
		resMap["msg"] = fmt.Sprintf("code=%d", code)
		resMap["count"] = 0
		return resMap
	}

	if bodyDst["data"] == nil {
		resMap["ret"] = false
		resMap["msg"] = "data不存在"
		resMap["count"] = 0
		return resMap
	}

	data := bodyDst["data"].(map[string]interface{})

	var count int
	if taskType == 11 {
		if data["fans"] == nil {
			resMap["ret"] = false
			resMap["msg"] = "fans不存在"
			resMap["count"] = 0
			return resMap
		}

		fansStr := data["fans"].(string)
		fans, ferr := strconv.Atoi(fansStr)
		if ferr != nil {
			resMap["ret"] = false
			resMap["msg"] = "fans转数字失败，" + fansStr
			resMap["count"] = 0
			return resMap
		}
		count = fans
	} else if taskType == 12 {
		if data["like"] == nil {
			resMap["ret"] = false
			resMap["msg"] = "fans不存在"
			resMap["count"] = 0
			return resMap
		}

		like := int(data["like"].(float64))
		count = like
	}

	resMap["ret"] = true
	resMap["msg"] = "获取成功"
	resMap["count"] = count
	return resMap
}
