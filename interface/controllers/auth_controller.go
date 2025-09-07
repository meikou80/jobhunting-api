package controllers

import (
	"net/http"

	"jobhunting-api/usecase/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthController 認証コントローラー
type AuthController struct {
	authService *services.AuthService
	validator   *validator.Validate
}

// NewAuthController 認証コントローラーのコンストラクタ
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		validator:   validator.New(),
	}
}

// Register ユーザー登録
// @Summary ユーザー登録
// @Description 新しいユーザーアカウントを作成します
// @Tags auth
// @Accept json
// @Produce json
// @Param user body services.RegisterRequest true "ユーザー登録情報"
// @Success 201 {object} services.AuthResponse "登録成功"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 409 {object} map[string]interface{} "メールアドレス重複"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req services.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	response, err := c.authService.Register(&req)
	if err != nil {
		if err.Error() == "email already exists" {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "Email already exists",
				"message": "An account with this email already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to register user",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    response,
	})
}

// Login ユーザーログイン
// @Summary ユーザーログイン
// @Description メールアドレスとパスワードでログインします
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body services.LoginRequest true "ログイン情報"
// @Success 200 {object} services.AuthResponse "ログイン成功"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 401 {object} map[string]interface{} "認証失敗"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req services.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	response, err := c.authService.Login(&req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"message": "Email or password is incorrect",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to login",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

// Logout ユーザーログアウト
// @Summary ユーザーログアウト
// @Description JWTトークンを無効化します（クライアント側でトークン削除）
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "ログアウト成功"
// @Security BearerAuth
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
		"note":    "Please remove the token from client side",
	})
}

// GetMe 現在のユーザー情報取得
// @Summary 現在のユーザー情報取得
// @Description 認証済みユーザーの情報を取得します
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "ユーザー情報"
// @Failure 401 {object} map[string]interface{} "認証が必要"
// @Security BearerAuth
// @Router /auth/me [get]
func (c *AuthController) GetMe(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userEmail, _ := ctx.Get("user_email")
	userRole, _ := ctx.Get("user_role")

	ctx.JSON(http.StatusOK, gin.H{
		"user": map[string]interface{}{
			"id":    userID,
			"email": userEmail,
			"role":  userRole,
		},
	})
}

// RefreshToken トークンリフレッシュ
// @Summary トークンリフレッシュ
// @Description 既存のトークンを使用して新しいトークンを発行します
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "新しいトークン"
// @Failure 401 {object} map[string]interface{} "認証が必要"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Security BearerAuth
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// 簡易実装: 現在のユーザー情報で新しいトークンを生成
	// 実際の実装では、データベースからユーザー情報を再取得することを推奨
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Token refresh not implemented yet",
		"user_id": userID,
		"note":    "Please login again to get a new token",
	})
}
