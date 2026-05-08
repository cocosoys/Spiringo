package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 中文：Validator 定义当前包使用的数据结构或接口。
// English: Validator defines a data structure or interface used by this package.
type Validator struct {
	// 中文：inner 保存当前结构中的配置或数据值。
	// English: inner stores a configuration or data value for this struct.
	inner *validator.Validate
}

// 中文：FieldLevel 定义当前包使用的数据结构或接口。
// English: FieldLevel defines a data structure or interface used by this package.
type FieldLevel = validator.FieldLevel

// 中文：Func 定义当前包使用的数据结构或接口。
// English: Func defines a data structure or interface used by this package.
type Func = validator.Func

// 中文：FieldError 定义当前包使用的数据结构或接口。
// English: FieldError defines a data structure or interface used by this package.
type FieldError struct {
	// 中文：Field 保存当前结构中的配置或数据值。
	// English: Field stores a configuration or data value for this struct.
	Field string `json:"field"`
	// 中文：Tag 保存当前结构中的配置或数据值。
	// English: Tag stores a configuration or data value for this struct.
	Tag string `json:"tag"`
	// 中文：Message 保存当前结构中的配置或数据值。
	// English: Message stores a configuration or data value for this struct.
	Message string `json:"message"`
}

// 中文：New 创建并返回对应组件实例。
// English: New creates and returns the corresponding component instance.
func New() *Validator {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "-" || name == "" {
			return field.Name
		}
		return name
	})
	return &Validator{inner: v}
}

// 中文：Struct 执行当前包中的对应流程。
// English: Struct executes the corresponding workflow in this package.
func (v *Validator) Struct(value any) error {
	return v.inner.Struct(value)
}

// 中文：Var 执行当前包中的对应流程。
// English: Var executes the corresponding workflow in this package.
func (v *Validator) Var(value any, tag string) error {
	return v.inner.Var(value, tag)
}

// 中文：RegisterValidation 执行当前包中的对应流程。
// English: RegisterValidation executes the corresponding workflow in this package.
func (v *Validator) RegisterValidation(tag string, fn Func) error {
	return v.inner.RegisterValidation(tag, fn)
}

// 中文：FieldErrors 执行当前包中的对应流程。
// English: FieldErrors executes the corresponding workflow in this package.
func (v *Validator) FieldErrors(err error) []FieldError {
	if err == nil {
		return nil
	}
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []FieldError{{Message: err.Error()}}
	}
	result := make([]FieldError, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		result = append(result, FieldError{
			Field:   fieldErr.Field(),
			Tag:     fieldErr.Tag(),
			Message: defaultMessage(fieldErr),
		})
	}
	return result
}

// 中文：BindJSON 执行当前包中的对应流程。
// English: BindJSON executes the corresponding workflow in this package.
func BindJSON(c *gin.Context, target any) error {
	if err := c.ShouldBindJSON(target); err != nil {
		return err
	}
	return New().Struct(target)
}

// 中文：defaultMessage 执行当前包中的对应流程。
// English: defaultMessage executes the corresponding workflow in this package.
func defaultMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s failed validation %s", err.Field(), err.Tag())
	}
}
