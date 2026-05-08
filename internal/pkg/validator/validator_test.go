package validator

import "testing"

// 中文：sampleRequest 定义当前包使用的数据结构或接口。
// English: sampleRequest defines a data structure or interface used by this package.
type sampleRequest struct {
	// 中文：Email 保存当前结构中的配置或数据值。
	// English: Email stores a configuration or data value for this struct.
	Email string `json:"email" validate:"required,email"`
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name" validate:"min=3"`
}

// 中文：TestValidatorUsesJSONFieldNames 验证相关行为符合预期。
// English: TestValidatorUsesJSONFieldNames verifies the related behavior.
func TestValidatorUsesJSONFieldNames(t *testing.T) {
	v := New()
	err := v.Struct(sampleRequest{Name: "golang"})
	if err == nil {
		t.Fatal("expected validation error")
	}

	fields := v.FieldErrors(err)
	if len(fields) != 1 {
		t.Fatalf("expected one error, got %#v", fields)
	}
	if fields[0].Field != "email" || fields[0].Tag != "required" {
		t.Fatalf("unexpected field error: %#v", fields[0])
	}
}

// 中文：TestRegisterValidation 验证相关行为符合预期。
// English: TestRegisterValidation verifies the related behavior.
func TestRegisterValidation(t *testing.T) {
	v := New()
	if err := v.RegisterValidation("even", func(fl FieldLevel) bool {
		return fl.Field().Int()%2 == 0
	}); err != nil {
		t.Fatal(err)
	}
	if err := v.Var(3, "even"); err == nil {
		t.Fatal("expected odd value to fail")
	}
}
