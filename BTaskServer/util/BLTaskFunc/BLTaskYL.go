package BLTaskFunc

import (
	"BTaskServer/util/Tools"
	"crypto/sha1"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

// 获取订单列表
func GetOrderListYL(page int, list_rows int, status int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.亿乐SUP")
	url := "/openapi/supplier/Order/Paging"

	AppID := viper.GetString("Ukey.AppID")
	AppSecret := viper.GetString("Ukey.AppSecret")
	AppTimestamp := strconv.Itoa(int(Tools.GetNowUnix(true)))
	AppTokenStr := fmt.Sprintf("%s%s%s%s", AppID, AppSecret, url, AppTimestamp)
	AppToken := getAppToken(AppTokenStr)

	rurl := ip + url
	headerMap := make(map[string]string)
	headerMap["Appid"] = AppID
	headerMap["AppTimestamp"] = AppTimestamp
	headerMap["AppToken"] = AppToken

	jsonMap := make(map[string]interface{})
	jsonMap["page"] = page
	jsonMap["list_rows"] = list_rows
	jsonMap["status"] = status

	bodyDst, err := Tools.PostJsonRequestTool(rurl, jsonMap, headerMap)

	if err != nil {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = "请求错误" + err.Error()
		return resMap
	}

	if bodyDst["code"] == nil {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = "code不存在"
		return resMap
	}

	code := int(bodyDst["code"].(float64))
	if code != 0 {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = fmt.Sprintf("code:%d,%s", code, bodyDst["message"].(string))
		return resMap
	}

	if bodyDst["data"] == nil {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = "data不存在"
		return resMap
	}

	data := bodyDst["data"]

	resMap["ret"] = true
	resMap["data"] = data
	resMap["msg"] = "获取成功"
	return resMap
}

// 获取播放订单列表
func GetBfOrderListYL(page int, list_rows int, goods_id int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.亿乐SUP")
	url := "/openapi/supplier/Order/Paging"

	AppID := viper.GetString("Ukey.AppID")
	AppSecret := viper.GetString("Ukey.AppSecret")
	AppTimestamp := strconv.Itoa(int(Tools.GetNowUnix(true)))
	AppTokenStr := fmt.Sprintf("%s%s%s%s", AppID, AppSecret, url, AppTimestamp)
	AppToken := getAppToken(AppTokenStr)

	rurl := ip + url
	headerMap := make(map[string]string)
	headerMap["Appid"] = AppID
	headerMap["AppTimestamp"] = AppTimestamp
	headerMap["AppToken"] = AppToken

	jsonMap := make(map[string]interface{})
	jsonMap["page"] = page
	jsonMap["list_rows"] = list_rows
	jsonMap["goods_id"] = goods_id

	bodyDst, err := Tools.PostJsonRequestTool(rurl, jsonMap, headerMap)

	if err != nil {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = "请求错误" + err.Error()
		return resMap
	}

	if bodyDst["code"] == nil {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = "code不存在"
		return resMap
	}

	code := int(bodyDst["code"].(float64))
	if code != 0 {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = fmt.Sprintf("code:%d,%s", code, bodyDst["message"].(string))
		return resMap
	}

	if bodyDst["data"] == nil {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = "data不存在"
		return resMap
	}

	data := bodyDst["data"]

	resMap["ret"] = true
	resMap["data"] = data
	resMap["msg"] = "获取成功"
	return resMap
}

// 修改订单状态
func UpdateOrderStatusYL(orderId int, oldStatus int, newStatus int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.亿乐SUP")
	url := "/openapi/supplier/Order/StatusHandle"

	AppID := viper.GetString("Ukey.AppID")
	AppSecret := viper.GetString("Ukey.AppSecret")
	AppTimestamp := strconv.Itoa(int(Tools.GetNowUnix(true)))
	AppTokenStr := fmt.Sprintf("%s%s%s%s", AppID, AppSecret, url, AppTimestamp)
	AppToken := getAppToken(AppTokenStr)

	rurl := ip + url
	headerMap := make(map[string]string)
	headerMap["Appid"] = AppID
	headerMap["AppTimestamp"] = AppTimestamp
	headerMap["AppToken"] = AppToken

	jsonMap := make(map[string]interface{})
	jsonMap["old_status"] = oldStatus
	jsonMap["new_status"] = newStatus
	jsonMap["id"] = orderId

	bodyDst, err := Tools.PostJsonRequestTool(rurl, jsonMap, headerMap)

	if err != nil {
		resMap["ret"] = false
		resMap["msg"] = "请求错误"
		return resMap
	}

	if bodyDst["code"] == nil {
		resMap["ret"] = false
		resMap["msg"] = "code不存在"
		return resMap
	}

	code := int(bodyDst["code"].(float64))
	if code != 0 {
		resMap["ret"] = false
		resMap["msg"] = fmt.Sprintf("code:%d,%s", code, bodyDst["message"].(string))
		return resMap
	}

	resMap["ret"] = true
	resMap["msg"] = "更新成功"
	return resMap
}

// 修改订单进度
func UpdateOrderScheduleYL(orderId int, startNum int, currentNum int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.亿乐SUP")
	url := "/openapi/supplier/Order/ScheduleHandle"

	AppID := viper.GetString("Ukey.AppID")
	AppSecret := viper.GetString("Ukey.AppSecret")
	AppTimestamp := strconv.Itoa(int(Tools.GetNowUnix(true)))
	AppTokenStr := fmt.Sprintf("%s%s%s%s", AppID, AppSecret, url, AppTimestamp)
	AppToken := getAppToken(AppTokenStr)

	rurl := ip + url
	headerMap := make(map[string]string)
	headerMap["Appid"] = AppID
	headerMap["AppTimestamp"] = AppTimestamp
	headerMap["AppToken"] = AppToken

	jsonMap := make(map[string]interface{})
	jsonMap["start_num"] = startNum
	jsonMap["current_num"] = currentNum
	jsonMap["id"] = orderId

	bodyDst, err := Tools.PostJsonRequestTool(rurl, jsonMap, headerMap)

	if err != nil {
		resMap["ret"] = false
		resMap["msg"] = "请求错误"
		return resMap
	}

	if bodyDst["code"] == nil {
		resMap["ret"] = false
		resMap["msg"] = "code不存在"
		return resMap
	}

	code := int(bodyDst["code"].(float64))
	if code != 0 {
		resMap["ret"] = false
		resMap["msg"] = fmt.Sprintf("code:%d,%s", code, bodyDst["message"].(string))
		return resMap
	}

	resMap["ret"] = true
	resMap["msg"] = "更新成功"
	return resMap
}

// 订单退款
func OrderrefundYL(orderId int, oldStatus int, newStatus int, remark string, refund_number int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.亿乐SUP")
	url := "/openapi/supplier/Order/StatusHandle"

	AppID := viper.GetString("Ukey.AppID")
	AppSecret := viper.GetString("Ukey.AppSecret")
	AppTimestamp := strconv.Itoa(int(Tools.GetNowUnix(true)))
	AppTokenStr := fmt.Sprintf("%s%s%s%s", AppID, AppSecret, url, AppTimestamp)
	AppToken := getAppToken(AppTokenStr)

	rurl := ip + url
	headerMap := make(map[string]string)
	headerMap["Appid"] = AppID
	headerMap["AppTimestamp"] = AppTimestamp
	headerMap["AppToken"] = AppToken

	jsonMap := make(map[string]interface{})
	jsonMap["old_status"] = oldStatus
	jsonMap["new_status"] = newStatus
	jsonMap["id"] = orderId
	jsonMap["remark"] = remark
	jsonMap["refund_number"] = refund_number

	bodyDst, err := Tools.PostJsonRequestTool(rurl, jsonMap, headerMap)

	if err != nil {
		resMap["ret"] = false
		resMap["msg"] = "请求错误"
		return resMap
	}

	if bodyDst["code"] == nil {
		resMap["ret"] = false
		resMap["msg"] = "code不存在"
		return resMap
	}

	code := int(bodyDst["code"].(float64))
	if code != 0 {
		resMap["ret"] = false
		resMap["msg"] = fmt.Sprintf("code:%d,%s", code, bodyDst["message"].(string))
		return resMap
	}

	resMap["ret"] = true
	resMap["msg"] = "更新成功"
	return resMap
}

func getAppToken(dataStr string) string {
	data := []byte(dataStr)
	sha := sha1.Sum(data)
	shaStr := fmt.Sprintf("%x", sha)
	return shaStr
}
