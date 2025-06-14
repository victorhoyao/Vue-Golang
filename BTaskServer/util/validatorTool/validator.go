package validatorTool

import (
	"BTaskServer/common"
	"BTaskServer/util/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 参数校验 json
func ValidatorJson[M any](c *gin.Context, paramsModel M) bool {
	if err := c.ShouldBindJSON(paramsModel); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非valiadator类型的错误直接返回
			response.RequsetBad(c, nil, errs.Error())
			return false
		}

		// validatorTool.ValidationErrors类型错误则进行翻译
		// 并使用removeTopStruct函数去除字段名中的结构体名称标识
		translatedErrors := common.RemoveTopStruct(errs.Translate(common.Trans))

		// Join all error messages into a single string
		var errorMessages []string
		for _, msg := range translatedErrors {
			errorMessages = append(errorMessages, msg)
		}
		finalErrorMessage := strings.Join(errorMessages, ", ")

		response.Fail(c, nil, finalErrorMessage)
		return false
	}
	//校验成功
	return true
}

// 参数校验 form-data
func ValidatorForm[M any](c *gin.Context, paramsModel M) bool {
	if err := c.ShouldBind(paramsModel); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非valiadator类型的错误直接返回
			response.RequsetBad(c, nil, errs.Error())
			return false
		}
		// validatorTool.ValidationErrors类型错误则进行翻译
		// 并使用removeTopStruct函数去除字段名中的结构体名称标识
		translatedErrors := common.RemoveTopStruct(errs.Translate(common.Trans))
		var errorMessages []string
		for _, msg := range translatedErrors {
			errorMessages = append(errorMessages, msg)
		}
		finalErrorMessage := strings.Join(errorMessages, ", ")
		response.Fail(c, nil, finalErrorMessage)
		return false
	}
	//校验成功
	return true
}

// 参数校验 query
func ValidatorQuery[M any](c *gin.Context, paramsModel M) bool {
	if err := c.ShouldBindQuery(paramsModel); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非valiadator类型的错误直接返回
			response.RequsetBad(c, nil, errs.Error())
			return false
		}
		// validatorTool.ValidationErrors类型错误则进行翻译
		// 并使用removeTopStruct函数去除字段名中的结构体名称标识
		translatedErrors := common.RemoveTopStruct(errs.Translate(common.Trans))
		var errorMessages []string
		for _, msg := range translatedErrors {
			errorMessages = append(errorMessages, msg)
		}
		finalErrorMessage := strings.Join(errorMessages, ", ")
		response.Fail(c, nil, finalErrorMessage)
		return false
	}
	//校验成功
	return true
}

// 参数校验 uri
func ValidatorUri(c *gin.Context, uriQuery string) (int, bool) {
	id, err := strconv.Atoi(c.Params.ByName(uriQuery))
	if err != nil {
		response.Fail(c, nil, "参数错误")
		return 0, false
	}
	return id, true
}
