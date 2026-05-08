package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// 中文：GenerateRSAKey 执行当前包中的对应流程。
// English: GenerateRSAKey executes the corresponding workflow in this package.
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	if bits < 2048 {
		return nil, fmt.Errorf("rsa key size must be at least 2048 bits")
	}
	return rsa.GenerateKey(rand.Reader, bits)
}

// 中文：EncodeRSAPrivateKeyPEM 执行当前包中的对应流程。
// English: EncodeRSAPrivateKeyPEM executes the corresponding workflow in this package.
func EncodeRSAPrivateKeyPEM(key *rsa.PrivateKey) []byte {
	if key == nil {
		return nil
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}

// 中文：EncodeRSAPublicKeyPEM 执行当前包中的对应流程。
// English: EncodeRSAPublicKeyPEM executes the corresponding workflow in this package.
func EncodeRSAPublicKeyPEM(key *rsa.PublicKey) ([]byte, error) {
	if key == nil {
		return nil, fmt.Errorf("rsa public key is required")
	}
	data, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: data}), nil
}

// 中文：ParseRSAPrivateKeyPEM 执行当前包中的对应流程。
// English: ParseRSAPrivateKeyPEM executes the corresponding workflow in this package.
func ParseRSAPrivateKeyPEM(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid private key pem")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("pem is not an rsa private key")
	}
	return key, nil
}

// 中文：ParseRSAPublicKeyPEM 执行当前包中的对应流程。
// English: ParseRSAPublicKeyPEM executes the corresponding workflow in this package.
func ParseRSAPublicKeyPEM(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid public key pem")
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		if cert, certErr := x509.ParseCertificate(block.Bytes); certErr == nil {
			parsed = cert.PublicKey
		} else {
			return nil, err
		}
	}
	key, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("pem is not an rsa public key")
	}
	return key, nil
}

// 中文：RSAEncrypt 执行当前包中的对应流程。
// English: RSAEncrypt executes the corresponding workflow in this package.
func RSAEncrypt(publicKey *rsa.PublicKey, plaintext []byte) (string, error) {
	if publicKey == nil {
		return "", fmt.Errorf("rsa public key is required")
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 中文：RSADecrypt 执行当前包中的对应流程。
// English: RSADecrypt executes the corresponding workflow in this package.
func RSADecrypt(privateKey *rsa.PrivateKey, encoded string) ([]byte, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("rsa private key is required")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
}

// 中文：RSASign 执行当前包中的对应流程。
// English: RSASign executes the corresponding workflow in this package.
func RSASign(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	if privateKey == nil {
		return "", fmt.Errorf("rsa private key is required")
	}
	sum := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, sum[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// 中文：RSAVerify 执行当前包中的对应流程。
// English: RSAVerify executes the corresponding workflow in this package.
func RSAVerify(publicKey *rsa.PublicKey, data []byte, encodedSignature string) error {
	if publicKey == nil {
		return fmt.Errorf("rsa public key is required")
	}
	signature, err := base64.StdEncoding.DecodeString(encodedSignature)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, sum[:], signature)
}
