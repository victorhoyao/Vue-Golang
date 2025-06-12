// 易客
package BLTaskFunc

import (
	"BTaskServer/util/Tools"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

// 获取订单列表
func GetOrderListYK(page int, list_rows int, status int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.易客")
	url := fmt.Sprintf("/api/supplier/order/v2/orders?page=%d&goodsSN=&orderSN=&state=%d&limit=%d", page, status, list_rows)

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

	bodyDst, err := Tools.GetRequestTool(rurl, headerMap)

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
	if code != 100 {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = fmt.Sprintf("code:%d,%s", code, bodyDst["message"].(string))
		return resMap
	}

	if bodyDst["result"] == nil {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = "data不存在"
		return resMap
	}

	data := bodyDst["result"]

	resMap["ret"] = true
	resMap["data"] = data
	resMap["msg"] = "获取成功"
	return resMap
}

// 修改订单状态
func UpdateOrderStatusYK(orderId int, oldStatus int, newStatus int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.易客")
	url := "/api/supplier/order/v2/order"

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
	jsonMap["orderSN"] = orderId
	jsonMap["state"] = newStatus
	jsonMap["remarks"] = ""

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
	if code != 100 {
		resMap["ret"] = false
		resMap["msg"] = fmt.Sprintf("code:%d,%s", code, bodyDst["message"].(string))
		return resMap
	}

	resMap["ret"] = true
	resMap["msg"] = "更新成功"
	return resMap
}

// 修改订单进度
func UpdateOrderScheduleYK(orderId int, startNum int, currentNum int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.易客")
	url := "/api/supplier/order/v2/progress"

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
	jsonMap["orderSN"] = orderId
	jsonMap["startNum"] = startNum
	jsonMap["currentNum"] = currentNum
	jsonMap["remarks"] = ""

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
	if code != 100 {
		resMap["ret"] = false
		resMap["msg"] = fmt.Sprintf("code:%d,%s", code, bodyDst["msg"].(string))
		return resMap
	}

	resMap["ret"] = true
	resMap["msg"] = "更新成功"
	return resMap
}

// 获取播放订单列表
func GetBfOrderListYK(page int, list_rows int, goods_id int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.易客")
	url := fmt.Sprintf("/api/supplier/order/v2/orders?page=%d&goodsSN=%d&orderSN=&state=&limit=%d", page, goods_id, list_rows)

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

	bodyDst, err := Tools.GetRequestTool(rurl, headerMap)

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
	if code != 100 {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = fmt.Sprintf("code:%d,%s", code, bodyDst["message"].(string))
		return resMap
	}

	if bodyDst["result"] == nil {
		resMap["ret"] = false
		resMap["data"] = nil
		resMap["msg"] = "data不存在"
		return resMap
	}

	data := bodyDst["result"]

	resMap["ret"] = true
	resMap["data"] = data
	resMap["msg"] = "获取成功"
	return resMap
}

// 订单退款
func OrderrefundYK(orderId int, oldStatus int, newStatus int, remark string, refund_number int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.易客")
	url := "/api/supplier/order/v2/refund"

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
	jsonMap["orderSN"] = orderId
	jsonMap["number"] = refund_number
	jsonMap["remarks"] = ""

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
	if code != 100 {
		resMap["ret"] = false
		resMap["msg"] = fmt.Sprintf("code:%d,%s", code, bodyDst["message"].(string))
		return resMap
	}

	resMap["ret"] = true
	resMap["msg"] = "更新成功"
	return resMap
}
