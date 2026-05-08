package types

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// 中文：TestAppError_WithMessage 验证相关行为符合预期。
// English: TestAppError_WithMessage verifies the related behavior.
func TestAppError_WithMessage(t *testing.T) {
	err := ErrBadRequest.WithMessage("custom message")
	if err.Code != ErrBadRequest.Code {
		t.Errorf("expected code %d, got %d", ErrBadRequest.Code, err.Code)
	}
	if err.Message != "custom message" {
		t.Errorf("expected 'custom message', got '%s'", err.Message)
	}
}

// 中文：TestAppError_WithMessagef 验证相关行为符合预期。
// English: TestAppError_WithMessagef verifies the related behavior.
func TestAppError_WithMessagef(t *testing.T) {
	err := ErrNotFound.WithMessagef("item %s not found", "abc")
	if err.Message != "item abc not found" {
		t.Errorf("expected 'item abc not found', got '%s'", err.Message)
	}
}

// 中文：TestAppError_Error 验证相关行为符合预期。
// English: TestAppError_Error verifies the related behavior.
func TestAppError_Error(t *testing.T) {
	err := ErrBadRequest.WithMessage("test")
	// Error() returns format [code] message
	expected := "[10001] test"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

// 中文：TestResponse_JSON 验证相关行为符合预期。
// English: TestResponse_JSON verifies the related behavior.
func TestResponse_JSON(t *testing.T) {
	resp := Response{Code: 0, Message: "ok", Data: map[string]string{"key": "value"}}
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	var parsed Response
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed.Code != 0 || parsed.Message != "ok" {
		t.Errorf("unexpected response: %+v", parsed)
	}
}

// 中文：TestOK 验证相关行为符合预期。
// English: TestOK verifies the related behavior.
func TestOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	OK(c, map[string]string{"hello": "world"})

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// 中文：TestFail 验证相关行为符合预期。
// English: TestFail verifies the related behavior.
func TestFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Fail(c, ErrBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// 中文：TestFail_UnknownError 验证相关行为符合预期。
// English: TestFail_UnknownError verifies the related behavior.
func TestFail_UnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Fail(c, http.ErrNoCookie)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// 中文：TestOKWithPage 验证相关行为符合预期。
// English: TestOKWithPage verifies the related behavior.
func TestOKWithPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	OKWithPage(c, []string{"a", "b"}, 2, 1, 10)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// JSON反序列化后Data是map[string]interface{}，需要二次序列化
	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	dataBytes, _ := json.Marshal(resp.Data)
	var data PageData
	json.Unmarshal(dataBytes, &data)
	if data.Total != 2 || data.Page != 1 || data.PageSize != 10 {
		t.Errorf("unexpected page data: %+v", data)
	}
}

// 中文：TestPaginationRequest_Defaults 验证相关行为符合预期。
// English: TestPaginationRequest_Defaults verifies the related behavior.
func TestPaginationRequest_Defaults(t *testing.T) {
	req := PaginationRequest{}
	if req.GetPage() != 1 {
		t.Errorf("expected default page 1, got %d", req.GetPage())
	}
	if req.GetPageSize() != 20 {
		t.Errorf("expected default page_size 20, got %d", req.GetPageSize())
	}
}

// 中文：TestPaginationRequest_Overflow 验证相关行为符合预期。
// English: TestPaginationRequest_Overflow verifies the related behavior.
func TestPaginationRequest_Overflow(t *testing.T) {
	req := PaginationRequest{Page: 0, PageSize: 200}
	if req.GetPage() != 1 {
		t.Errorf("expected page 1, got %d", req.GetPage())
	}
	if req.GetPageSize() != 100 {
		t.Errorf("expected page_size capped at 100, got %d", req.GetPageSize())
	}
}
