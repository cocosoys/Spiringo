package snowflake

import (
	"testing"
)

// 中文：TestSnowflake_Generate 验证相关行为符合预期。
// English: TestSnowflake_Generate verifies the related behavior.
func TestSnowflake_Generate(t *testing.T) {
	sf, err := New(1)
	if err != nil {
		t.Fatal(err)
	}

	id1 := sf.Generate()
	id2 := sf.Generate()

	if id1 == 0 {
		t.Error("ID should not be zero")
	}
	if id1 == id2 {
		t.Error("consecutive IDs should be different")
	}
	if id2 < id1 {
		t.Error("IDs should be monotonically increasing")
	}
}

// 中文：TestSnowflake_Uniqueness 验证相关行为符合预期。
// English: TestSnowflake_Uniqueness verifies the related behavior.
func TestSnowflake_Uniqueness(t *testing.T) {
	sf, _ := New(1)

	ids := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id := sf.Generate()
		if ids[id] {
			t.Errorf("duplicate ID generated: %d", id)
		}
		ids[id] = true
	}
}

// 中文：TestSnowflake_InvalidMachineID 验证相关行为符合预期。
// English: TestSnowflake_InvalidMachineID verifies the related behavior.
func TestSnowflake_InvalidMachineID(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Error("expected error for negative machine ID")
	}

	_, err = New(1024)
	if err == nil {
		t.Error("expected error for machine ID > 1023")
	}
}

// 中文：TestSnowflake_DifferentMachines 验证相关行为符合预期。
// English: TestSnowflake_DifferentMachines verifies the related behavior.
func TestSnowflake_DifferentMachines(t *testing.T) {
	sf1, _ := New(1)
	sf2, _ := New(2)

	id1 := sf1.Generate()
	id2 := sf2.Generate()

	// Same timestamp+sequence, different machine should produce different IDs
	if id1 == id2 {
		t.Error("different machines should produce different IDs")
	}
}

// 中文：TestSnowflake_GenerateString 验证相关行为符合预期。
// English: TestSnowflake_GenerateString verifies the related behavior.
func TestSnowflake_GenerateString(t *testing.T) {
	sf, _ := New(1)

	str := sf.GenerateString()
	if str == "" {
		t.Error("string ID should not be empty")
	}
	if str == "0" {
		t.Error("string ID should not be '0'")
	}
}
