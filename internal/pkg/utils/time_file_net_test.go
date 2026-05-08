package utils

import (
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

// 中文：TestTimeHelpers 验证相关行为符合预期。
// English: TestTimeHelpers verifies the related behavior.
func TestTimeHelpers(t *testing.T) {
	parsed, err := ParseTime("2026-05-08 12:30:00")
	if err != nil {
		t.Fatal(err)
	}
	if FormatTime(parsed) != "2026-05-08 12:30:00" {
		t.Fatalf("unexpected formatted time: %s", FormatTime(parsed))
	}
	if StartOfDay(parsed).Hour() != 0 || EndOfDay(parsed).Hour() != 23 {
		t.Fatalf("unexpected day bounds: %v %v", StartOfDay(parsed), EndOfDay(parsed))
	}
}

// 中文：TestJSONFileHelpers 验证相关行为符合预期。
// English: TestJSONFileHelpers verifies the related behavior.
func TestJSONFileHelpers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.json")
	value := map[string]string{"name": "spiringo"}
	if err := WriteJSON(path, value); err != nil {
		t.Fatal(err)
	}
	if !FileExists(path) || !DirExists(filepath.Dir(path)) {
		t.Fatal("expected json file and parent directory to exist")
	}

	var out map[string]string
	if err := ReadJSON(path, &out); err != nil {
		t.Fatal(err)
	}
	if out["name"] != "spiringo" {
		t.Fatalf("unexpected json: %#v", out)
	}
	if err := RemoveIfExists(path); err != nil {
		t.Fatal(err)
	}
	if err := RemoveIfExists(path); err != nil {
		t.Fatal(err)
	}
}

// 中文：TestClientIP 验证相关行为符合预期。
// English: TestClientIP verifies the related behavior.
func TestClientIP(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Forwarded-For", "203.0.113.10, 10.0.0.1")
	req.RemoteAddr = "127.0.0.1:1234"

	if got := ClientIP(req); got != "203.0.113.10" {
		t.Fatalf("ClientIP = %q", got)
	}
	if !IsPrivateIP("127.0.0.1") {
		t.Fatal("loopback should be treated as private")
	}
}

// 中文：TestNowUnixMilli 验证相关行为符合预期。
// English: TestNowUnixMilli verifies the related behavior.
func TestNowUnixMilli(t *testing.T) {
	before := time.Now().Add(-time.Second).UnixMilli()
	now := NowUnixMilli()
	after := time.Now().Add(time.Second).UnixMilli()
	if now < before || now > after {
		t.Fatalf("NowUnixMilli outside expected range: %d", now)
	}
}
