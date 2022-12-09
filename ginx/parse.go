package ginx

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func parseErrorTag(field, tag string) string {
	if field == "" {
		// 排除未取到field名的错误
		return ""
	}

	switch tag {
	case "required":
		return fmt.Sprintf("参数 %s 是必填参数。", field)
	case "email":
		return "Invalid email。"
	}
	return fmt.Sprintf("参数 %s 不合法。", field)
}

func parseParamErrorDetails(t reflect.Type, err error, tag string) string {
	// https://github.com/gin-gonic/gin/issues/2334
	// https://stackoverflow.com/questions/70069834/return-custom-error-message-from-struct-tag-validation
	details := "Bad Request。具体是："
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field, _ := t.Elem().FieldByName(fe.Field())
			fieldName, _ := field.Tag.Lookup(tag)
			details += parseErrorTag(fieldName, fe.Tag())
		}
	} else {
		details += err.Error()
	}

	return details
}

func ParseJsonRequest(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		// 参数不合法
		desc := parseParamErrorDetails(reflect.TypeOf(obj), err, "json")
		return errors.New(desc)
	}

	return nil
}

func ParseQueryRequest(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		// 参数不合法
		desc := parseParamErrorDetails(reflect.TypeOf(obj), err, "form")
		return errors.New(desc)
	}

	return nil
}

func ParseFormRequest(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		// 参数不合法
		desc := parseParamErrorDetails(reflect.TypeOf(obj), err, "form")
		return errors.New(desc)
	}

	return nil
}

func ParseRequest(c *gin.Context, obj interface{}) error {
	// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）。
	// 如果是 `POST` 请求，首先检查 `content-type` 是否为 `JSON` 或 `XML`，然后再使用 `Form`（`form-data`）。
	// 查看更多：https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L88
	if err := c.ShouldBind(obj); err != nil {
		// 参数不合法
		desc := parseParamErrorDetails(reflect.TypeOf(obj), err, "form")
		return errors.New(desc)
	}

	return nil
}
