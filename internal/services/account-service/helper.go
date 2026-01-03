package accountservice

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
)

// maskAPIKey masks sensitive API key information for logging purposes
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 10 {
		return apiKey
	}
	return apiKey[:10] + "..."
}

// generateAPIKey generates a secure random API key
func generateAPIKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return "api_" + string(b)
}

// generateMerchantID generates a unique merchant ID with CASM- prefix
func generateMerchantID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 12)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return "CASM-" + string(b)
}

// generateSecretKey generates a secure random secret key (currently unused)
func generateSecretKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 64)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return "sk_" + string(b)
}

// hashAPIKey uses HMAC-SHA256 for deterministic hashing of API keys for storage and comparison
func hashAPIKey(apiKey, hashKey string) string {
	h := hmac.New(sha256.New, []byte(hashKey))
	h.Write([]byte(apiKey))
	hashed := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashed)
}

// encryptAPIKey uses AES-GCM for encrypting the plain API key for secure storage
func encryptAPIKey(apiKey, hashKey string) string {
	block, err := aes.NewCipher([]byte(hashKey)[:32])
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(apiKey), nil)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// decryptAPIKey decrypts an AES-GCM encrypted API key
func decryptAPIKey(encryptedKey, hashKey string) (string, error) {
	block, err := aes.NewCipher([]byte(hashKey)[:32])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
