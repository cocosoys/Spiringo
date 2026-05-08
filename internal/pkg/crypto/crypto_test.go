package crypto

import (
	"testing"
)

// 中文：TestAES_EncryptDecrypt 验证相关行为符合预期。
// English: TestAES_EncryptDecrypt verifies the related behavior.
func TestAES_EncryptDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef") // 32 bytes
	plaintext := []byte("hello spiringo")

	ciphertext, err := AESEncrypt(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := AESDecrypt(key, ciphertext)
	if err != nil {
		t.Fatal(err)
	}

	if string(decrypted) != "hello spiringo" {
		t.Errorf("expected 'hello spiringo', got '%s'", string(decrypted))
	}
}

// 中文：TestAES_InvalidKey 验证相关行为符合预期。
// English: TestAES_InvalidKey verifies the related behavior.
func TestAES_InvalidKey(t *testing.T) {
	_, err := AESEncrypt([]byte("short"), []byte("test"))
	if err == nil {
		t.Error("expected error for invalid key size")
	}
}

// 中文：TestSHA256 验证相关行为符合预期。
// English: TestSHA256 verifies the related behavior.
func TestSHA256(t *testing.T) {
	hash := SHA256([]byte("hello"))
	if len(hash) != 64 {
		t.Errorf("expected 64 char hex hash, got %d", len(hash))
	}

	// Same input should produce same hash
	hash2 := SHA256([]byte("hello"))
	if hash != hash2 {
		t.Error("SHA256 should be deterministic")
	}

	// Different input should produce different hash
	hash3 := SHA256([]byte("world"))
	if hash == hash3 {
		t.Error("different inputs should produce different hashes")
	}
}

// 中文：TestSHA256String 验证相关行为符合预期。
// English: TestSHA256String verifies the related behavior.
func TestSHA256String(t *testing.T) {
	hash := SHA256String("hello")
	if len(hash) != 64 {
		t.Errorf("expected 64 char hex hash, got %d", len(hash))
	}

	// Should match SHA256([]byte("hello"))
	expected := SHA256([]byte("hello"))
	if hash != expected {
		t.Error("SHA256String should match SHA256 for same input")
	}
}

// 中文：TestJWT_TokenGenerationAndParsing 验证相关行为符合预期。
// English: TestJWT_TokenGenerationAndParsing verifies the related behavior.
func TestJWT_TokenGenerationAndParsing(t *testing.T) {
	cfg := JWTConfig{
		Secret:        "test-secret-key-for-jwt",
		AccessExpire:  3600e9,  // 1 hour in nanoseconds
		RefreshExpire: 86400e9, // 1 day in nanoseconds
		Issuer:        "spiringo-test",
	}

	accessToken, refreshToken, err := GenerateToken(cfg, "user-123", "testuser", "tenant-1")
	if err != nil {
		t.Fatal(err)
	}

	if accessToken == "" {
		t.Error("access token should not be empty")
	}
	if refreshToken == "" {
		t.Error("refresh token should not be empty")
	}

	// Parse access token
	claims, err := ParseToken(cfg.Secret, accessToken)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("expected user_id 'user-123', got '%s'", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", claims.Username)
	}
	if claims.TenantID != "tenant-1" {
		t.Errorf("expected tenant_id 'tenant-1', got '%s'", claims.TenantID)
	}

	// Parse refresh token
	refreshClaims, err := ParseToken(cfg.Secret, refreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if refreshClaims.UserID != "user-123" {
		t.Errorf("expected user_id 'user-123', got '%s'", refreshClaims.UserID)
	}
}

// 中文：TestJWT_InvalidToken 验证相关行为符合预期。
// English: TestJWT_InvalidToken verifies the related behavior.
func TestJWT_InvalidToken(t *testing.T) {
	_, err := ParseToken("secret", "invalid.token.string")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

// 中文：TestRSAEncryptDecryptAndSign 验证相关行为符合预期。
// English: TestRSAEncryptDecryptAndSign verifies the related behavior.
func TestRSAEncryptDecryptAndSign(t *testing.T) {
	key, err := GenerateRSAKey(2048)
	if err != nil {
		t.Fatal(err)
	}

	privatePEM := EncodeRSAPrivateKeyPEM(key)
	parsedPrivate, err := ParseRSAPrivateKeyPEM(privatePEM)
	if err != nil {
		t.Fatal(err)
	}
	publicPEM, err := EncodeRSAPublicKeyPEM(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	parsedPublic, err := ParseRSAPublicKeyPEM(publicPEM)
	if err != nil {
		t.Fatal(err)
	}

	encrypted, err := RSAEncrypt(parsedPublic, []byte("hello rsa"))
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := RSADecrypt(parsedPrivate, encrypted)
	if err != nil {
		t.Fatal(err)
	}
	if string(decrypted) != "hello rsa" {
		t.Fatalf("decrypted = %q", decrypted)
	}

	signature, err := RSASign(parsedPrivate, []byte("payload"))
	if err != nil {
		t.Fatal(err)
	}
	if err := RSAVerify(parsedPublic, []byte("payload"), signature); err != nil {
		t.Fatal(err)
	}
	if err := RSAVerify(parsedPublic, []byte("tampered"), signature); err == nil {
		t.Fatal("expected signature verification to fail for tampered payload")
	}
}

// 中文：TestGenerateRSAKeyRejectsSmallKeys 验证相关行为符合预期。
// English: TestGenerateRSAKeyRejectsSmallKeys verifies the related behavior.
func TestGenerateRSAKeyRejectsSmallKeys(t *testing.T) {
	if _, err := GenerateRSAKey(1024); err == nil {
		t.Fatal("expected small rsa key to fail")
	}
}
