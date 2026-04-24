package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// 注册自定义验证规则
	_ = validate.RegisterValidation("fanhao", validateFanhao)
	_ = validate.RegisterValidation("username", validateUsername)
}

// Validate 验证结构体
func Validate(data interface{}) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	// 格式化错误信息
	var errs validator.ValidationErrors
	errors.As(err, &errs)
	var messages []string
	for _, e := range errs {
		messages = append(messages, formatError(e))
	}

	return fmt.Errorf(strings.Join(messages, "; "))
}

// formatError 格式化错误信息
func formatError(e validator.FieldError) string {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s不能为空", field)
	case "min":
		return fmt.Sprintf("%s长度不能小于%s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s长度不能大于%s", field, e.Param())
	case "email":
		return fmt.Sprintf("%s格式不正确", field)
	case "fanhao":
		return "番号格式不正确"
	case "username":
		return "用户名只能包含字母、数字和下划线"
	default:
		return fmt.Sprintf("%s验证失败", field)
	}
}

// validateFanhao 验证番号格式
func validateFanhao(fl validator.FieldLevel) bool {
	fanhao := fl.Field().String()
	// 番号格式：字母-数字，例如 SSIS-123
	matched, _ := regexp.MatchString(`^[A-Z]+-\d+$`, strings.ToUpper(fanhao))
	return matched
}

// validateUsername 验证用户名格式
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// 只允许字母、数字、下划线
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched
}
