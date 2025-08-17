package crypto

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims Supabase JWTトークンのクレーム
type JWTClaims struct {
	Sub   string `json:"sub"`   // ユーザーID (UUID)
	Email string `json:"email"` // ユーザーEmail
	Role  string `json:"role"`  // ユーザーロール
	jwt.RegisteredClaims
}

// JWTValidator JWT検証の責任を持つ構造体
type JWTValidator struct {
	jwtSecret string
}

// NewJWTValidator JWT検証インスタンスを生成
// 設計思想: 依存関係を明確化し、テスタビリティを向上
func NewJWTValidator() *JWTValidator {
	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
	if jwtSecret == "" {
		// 開発環境でのデフォルト値（本番では必須）
		jwtSecret = "your-256-bit-secret"
	}

	return &JWTValidator{
		jwtSecret: jwtSecret,
	}
}

// ValidateToken JWT トークンを検証し、ユーザー情報を抽出
// 設計思想:
// - 単一責任: JWT検証のみに集中
// - エラー詳細化: 問題特定を容易にする
// - 型安全性: 構造化されたクレーム情報を返却
func (v *JWTValidator) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Bearerプリフィックスを除去
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// JWT パース & 検証
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// アルゴリズム検証（セキュリティ対策）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(v.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// クレーム検証
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// 必須フィールド検証
	if claims.Sub == "" {
		return nil, errors.New("missing user ID in token")
	}

	return claims, nil
}

// ExtractUserID トークンからユーザーIDのみを抽出（簡易版）
// 設計思想: よく使われる操作を簡略化
func (v *JWTValidator) ExtractUserID(tokenString string) (string, error) {
	claims, err := v.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Sub, nil
}
