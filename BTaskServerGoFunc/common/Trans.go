package common

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

/*
参数校验翻译
*/
var Trans ut.Translator

// 初始化翻译器
func InintTrans() bool {
	locale := "zh"
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		Trans, ok = uni.GetTranslator("zh")
		if !ok {
			//return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
			fmt.Println("参数校验器启动失败")
			return false
		}

		// 注册翻译器
		switch locale {
		case "en":
			_ = enTranslations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			_ = zhTranslations.RegisterDefaultTranslations(v, Trans)
		default:
			_ = enTranslations.RegisterDefaultTranslations(v, Trans)
		}

		//{
		//	// 手机号自定义验证注册
		//	_ = v.RegisterValidation("telephone", validate.ValidateMobile)
		//	_ = v.RegisterTranslation("telephone", Trans, registerTranslator("telephone", "{0}格式不正确"), translate)
		//}

		//// Task自定义参数校验
		//{
		//	_ = v.RegisterValidation("isExceedNow", validate.ValidateTask)
		//	_ = v.RegisterTranslation("targetTime", Trans, registerTranslator("targetTime", "{0}目标时间不能比此刻更早！"), translate)
		//}
	}

	fmt.Println("参数校验器启动成功")
	return true
}

// 工具函数
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// registerTranslator 为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}
