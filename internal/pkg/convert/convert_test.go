package convert

import "testing"

// 中文：TestScalarConversions 验证相关行为符合预期。
// English: TestScalarConversions verifies the related behavior.
func TestScalarConversions(t *testing.T) {
	i, err := Int("42")
	if err != nil || i != 42 {
		t.Fatalf("Int = %d, %v", i, err)
	}
	f, err := Float64("3.5")
	if err != nil || f != 3.5 {
		t.Fatalf("Float64 = %v, %v", f, err)
	}
	b, err := Bool("true")
	if err != nil || !b {
		t.Fatalf("Bool = %v, %v", b, err)
	}
}

// 中文：TestStringsConversion 验证相关行为符合预期。
// English: TestStringsConversion verifies the related behavior.
func TestStringsConversion(t *testing.T) {
	values, err := Strings("a, b,,c")
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 3 || values[1] != "b" {
		t.Fatalf("unexpected strings: %#v", values)
	}
}

// 中文：TestJSONMap 验证相关行为符合预期。
// English: TestJSONMap verifies the related behavior.
func TestJSONMap(t *testing.T) {
	values, err := JSONMap([]byte(`{"name":"spiringo"}`))
	if err != nil {
		t.Fatal(err)
	}
	if values["name"] != "spiringo" {
		t.Fatalf("unexpected map: %#v", values)
	}
}
