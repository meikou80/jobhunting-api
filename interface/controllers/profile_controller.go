package controllers

import (
	"net/http"

	"jobhunting-api/domain/models"
	"jobhunting-api/infra/middleware"
	"jobhunting-api/usecase/services"

	"github.com/gin-gonic/gin"
)

// ProfileController ユーザープロフィール管理のコントローラー
type ProfileController struct {
	profileService *services.ProfileService
}

// NewProfileController コンストラクタ
func NewProfileController(profileService *services.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
	}
}

// GetProfile ユーザープロフィール取得
// GET /api/v1/profile
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
