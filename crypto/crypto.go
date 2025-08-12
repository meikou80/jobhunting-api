package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost bcryptのデフォルトコスト
	DefaultCost = 12
)

// HashPassword パスワードをハッシュ化
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	return string(bytes), err
}

// CheckPassword パスワードとハッシュを比較
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateFromPassword 指定したコストでパスワードをハッシュ化
func GenerateFromPassword(password string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = DefaultCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}
