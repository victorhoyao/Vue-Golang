// 易客
package BLTaskFunc

import (
	"BTaskServer/util/Tools"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

// 获取订单列表
func GetOrderListYK2(page int, list_rows int, status int) map[string]interface{} {
	resMap := make(map[string]interface{})

	ip := viper.GetString("平台ip.易客2")
	url := fmt.Sprintf("/api/supplier/order/v2/orders?page=%d&goodsSN=&orderSN=&state=%d&limit=%d", page, status, list_rows)

	AppID := viper.GetString("Ukey2.AppID")
	AppSecret := viper.GetString("Ukey2.AppSecret")
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
