package crypto

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 中文：JWTConfig 定义当前包使用的数据结构或接口。
// English: JWTConfig defines a data structure or interface used by this package.
// JWTConfig JWT配置
type JWTConfig struct {
	// 中文：Secret 保存当前结构中的配置或数据值。
	// English: Secret stores a configuration or data value for this struct.
	Secret string
	// 中文：AccessExpire 保存当前结构中的配置或数据值。
	// English: AccessExpire stores a configuration or data value for this struct.
	AccessExpire time.Duration
	// 中文：RefreshExpire 保存当前结构中的配置或数据值。
	// English: RefreshExpire stores a configuration or data value for this struct.
	RefreshExpire time.Duration
	// 中文：Issuer 保存当前结构中的配置或数据值。
	// English: Issuer stores a configuration or data value for this struct.
	Issuer string
}

// 中文：JWTClaims 定义当前包使用的数据结构或接口。
// English: JWTClaims defines a data structure or interface used by this package.
// JWTClaims JWT声明
type JWTClaims struct {
	// 中文：UserID 保存当前结构中的配置或数据值。
	// English: UserID stores a configuration or data value for this struct.
	UserID string `json:"user_id"`
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `json:"username"`
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `json:"tenant_id,omitempty"`
	// 中文：jwt.RegisteredClaims 嵌入复用该类型提供的能力。
	// English: jwt.RegisteredClaims embeds reusable behavior from that type.
	jwt.RegisteredClaims
}

// 中文：GenerateToken 执行当前包中的对应流程。
// English: GenerateToken executes the corresponding workflow in this package.
// GenerateToken 生成JWT Token
func GenerateToken(cfg JWTConfig, userID, username, tenantID string) (accessToken, refreshToken string, err error) {
	now := time.Now()

	accessClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.AccessExpire)),
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.RefreshExpire)),
		},
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", fmt.Errorf("sign refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// 中文：ParseToken 执行当前包中的对应流程。
// English: ParseToken executes the corresponding workflow in this package.
// ParseToken 解析JWT Token
func ParseToken(secret, tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
