package controllers

import (
	"net/http"

	"jobhunting-api/domain/models"
	"jobhunting-api/infra/middleware"
	"jobhunting-api/usecase/services"

	"github.com/gin-gonic/gin"
)

// ProfileController ユーザープロフィール管理のコントローラー
// 設計思想:
// - 薄いコントローラー: ビジネスロジックはServiceに委譲
// - HTTP層の責任: リクエスト/レスポンス変換、認証情報取得
// - エラーハンドリング: HTTPステータスコードとJSON形式の一貫性
type ProfileController struct {
	profileService *services.ProfileService
}

// NewProfileController コンストラクタ
// 設計思想: 依存性注入によりテスタビリティを向上
func NewProfileController(profileService *services.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
	}
}

// GetProfile ユーザープロフィール取得
// GET /api/v1/profile
// 設計思想:
// - 認証必須: JWTからユーザーIDを取得
// - 存在チェック: プロフィール未作成の場合は404
// - セキュリティ: 自分のプロフィールのみ取得可能
func (pc *ProfileController) GetProfile(c *gin.Context) {
	// 認証済みユーザーIDを取得
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_required",
			"message": "User authentication is required",
		})
		return
	}

	// プロフィール取得
	profile, err := pc.profileService.GetProfile(userID)
	if err != nil {
		// サービス層でのエラー種別に応じたHTTPステータス決定
		switch err.Error() {
		case "profile not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "profile_not_found",
				"message": "User profile not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Failed to retrieve profile",
				"detail":  err.Error(),
			})
		}
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}

// CreateProfile ユーザープロフィール作成
// POST /api/v1/profile
// 設計思想:
// - 初回登録: Supabase認証後の初回プロフィール作成
// - バリデーション: 必須フィールドと形式チェック
// - 冪等性考慮: 既存プロフィールがある場合の処理
func (pc *ProfileController) CreateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_required",
			"message": "User authentication is required",
		})
		return
	}

	// リクエストボディの解析
	var req models.UserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"detail":  err.Error(),
		})
		return
	}

	// バリデーションはGinのstructタグで自動実行される

	// プロフィール作成
	profile, err := pc.profileService.CreateProfile(userID, &req)
	if err != nil {
		switch err.Error() {
		case "profile already exists":
			c.JSON(http.StatusConflict, gin.H{
				"error":   "profile_exists",
				"message": "User profile already exists",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Failed to create profile",
				"detail":  err.Error(),
			})
		}
		return
	}

	// 作成成功
	c.JSON(http.StatusCreated, gin.H{
		"data":    profile,
		"message": "Profile created successfully",
	})
}

// UpdateProfile ユーザープロフィール更新
// PUT /api/v1/profile
// 設計思想:
// - 部分更新対応: 提供されたフィールドのみを更新
// - 楽観的ロック: 更新競合の検出
// - 入力検証: 更新データの妥当性確認
func (pc *ProfileController) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_required",
			"message": "User authentication is required",
		})
		return
	}

	var req models.UserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"detail":  err.Error(),
		})
		return
	}

	// 更新実行
	profile, err := pc.profileService.UpdateProfile(userID, &req)
	if err != nil {
		switch err.Error() {
		case "profile not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "profile_not_found",
				"message": "User profile not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Failed to update profile",
				"detail":  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    profile,
		"message": "Profile updated successfully",
	})
}

// DeleteProfile ユーザープロフィール削除（論理削除）
// DELETE /api/v1/profile
// 設計思想:
// - 論理削除: データの完全削除は避け、deleted_atで管理
// - 関連データ考慮: 求人・応募データとの整合性
// - 復旧可能性: 誤削除時の復旧を考慮
func (pc *ProfileController) DeleteProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_required",
			"message": "User authentication is required",
		})
		return
	}

	err := pc.profileService.DeleteProfile(userID)
	if err != nil {
		switch err.Error() {
		case "profile not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "profile_not_found",
				"message": "User profile not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Failed to delete profile",
				"detail":  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile deleted successfully",
	})
}

// GetUserSettings ユーザー設定取得
// GET /api/v1/profile/settings
// 設計思想: プロフィール情報と設定情報の分離
func (pc *ProfileController) GetUserSettings(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_required",
			"message": "User authentication is required",
		})
		return
	}

	settings, err := pc.profileService.GetUserSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to retrieve settings",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": settings,
	})
}

// UpdateUserSettings ユーザー設定更新
// PUT /api/v1/profile/settings
func (pc *ProfileController) UpdateUserSettings(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_required",
			"message": "User authentication is required",
		})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"detail":  err.Error(),
		})
		return
	}

	settings, err := pc.profileService.UpdateUserSettings(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to update settings",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    settings,
		"message": "Settings updated successfully",
	})
}
