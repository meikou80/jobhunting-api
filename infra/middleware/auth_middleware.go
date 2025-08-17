package middleware

import (
	"net/http"
	"strings"

	"jobhunting-api/crypto"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT認証ミドルウェア
// 設計思想: 認証処理を横断的関心事として分離し、各エンドポイントで再利用可能にする
type AuthMiddleware struct {
	jwtValidator *crypto.JWTValidator
}

// NewAuthMiddleware 認証ミドルウェアのコンストラクタ
// 設計思想: 依存性注入により、JWTValidator の実装を外部から注入可能
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		jwtValidator: crypto.NewJWTValidator(),
	}
}

// RequireAuth 認証必須のミドルウェア
// 設計思想:
// - 明示的な認証要求: エンドポイントが認証を必要とすることを明確化
// - 早期リターン: 認証失敗時は即座にエラーレスポンスを返却
// - コンテキスト設定: 認証成功時はユーザー情報をコンテキストに格納
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization ヘッダーからトークン抽出
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "authorization_required",
				"message": "Authorization header is required",
			})
			c.Abort() // 後続処理を停止
			return
		}

		// Bearer トークン形式の検証
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token_format",
				"message": "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// JWT トークン検証
		claims, err := m.jwtValidator.ValidateToken(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token",
				"message": "Invalid or expired token",
				"detail":  err.Error(), // 開発時のデバッグ用（本番では削除推奨）
			})
			c.Abort()
			return
		}

		// ユーザー情報をコンテキストに設定（後続処理で利用可能）
		// 設計思想: 型安全なコンテキストキーでデータ共有
		c.Set("user_id", claims.Sub)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_claims", claims) // 完全なクレーム情報も保存

		// 後続処理に進む
		c.Next()
	}
}

// OptionalAuth 認証任意のミドルウェア
// 設計思想:
// - 柔軟性: 認証されていてもいなくても処理を継続
// - コンテキスト一貫性: 認証済みの場合は RequireAuth と同じ形式でデータ設定
// - ログイン状態判定: 後続処理でログイン状態を判定可能
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// トークンが提供されていない場合は未認証として継続
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		// トークン検証（エラーの場合は未認証として継続）
		claims, err := m.jwtValidator.ValidateToken(authHeader)
		if err != nil {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		// 認証成功：ユーザー情報を設定
		c.Set("authenticated", true)
		c.Set("user_id", claims.Sub)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_claims", claims)

		c.Next()
	}
}

// GetUserID コンテキストからユーザーIDを取得するヘルパー関数
// 設計思想:
// - DRY原則: 共通処理を関数化
// - 型安全性: string型での取得を保証
// - エラーハンドリング: 取得失敗時の適切な処理
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

// GetUserClaims コンテキストから完全なユーザークレームを取得
// 設計思想: より詳細なユーザー情報が必要な場合の拡張性を提供
func GetUserClaims(c *gin.Context) (*crypto.JWTClaims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*crypto.JWTClaims)
	return userClaims, ok
}

// IsAuthenticated 認証状態の確認ヘルパー
// 設計思想: OptionalAuth使用時の認証状態判定を簡略化
func IsAuthenticated(c *gin.Context) bool {
	authenticated, exists := c.Get("authenticated")
	if !exists {
		// RequireAuth 使用時は authenticated フラグが設定されないため、
		// user_id の存在で判定
		_, hasUserID := c.Get("user_id")
		return hasUserID
	}

	auth, ok := authenticated.(bool)
	return ok && auth
}
