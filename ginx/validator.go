package ginx

import (
	"errors"
	"reflect"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

// use it:
//
//	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
//		RegisterValidatorTranslations(v)
//	}
//
// 注册 Validator Translations
func RegisterValidatorTranslations(v *validator.Validate) {
	english := en.New()
	chinese := zh.New()
	uni := ut.New(english, chinese, english)
	trans, _ = uni.GetTranslator("zh")
	zh_translations.RegisterDefaultTranslations(v, trans)
	// register tag for better prompt
	v.RegisterTagNameFunc(func(filed reflect.StructField) string {
		name := filed.Tag.Get("label")
		if name == "" {
			name = filed.Tag.Get("json")
		}
		if name == "" {
			name = filed.Tag.Get("form")
		}
		if name == "" {
			name = filed.Tag.Get("xml")
		}
		if name == "" {
			name = filed.Name
		}
		return name
	})
}

// Translate 翻译错误信息
func TranslateValidatorErrors(err error) []string {
	var results []string
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			results = append(results, fe.Translate(trans))
		}
	} else {
		results = append(results, err.Error())
	}
	return results
}

func TranslateValidatorErrorsString(err error) string {
	var details string
	results := TranslateValidatorErrors(err)
	for _, r := range results {
		if details != "" {
			details += "；" + r
		} else {
			details += r
		}
	}
	return details
}
