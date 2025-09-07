package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"jobhunting-api/crypto"
	"jobhunting-api/domain/models"
	"jobhunting-api/domain/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthService 認証サービス
type AuthService struct {
	authRepo    repositories.AuthRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

// NewAuthService 認証サービスのコンストラクタ
func NewAuthService(authRepo repositories.AuthRepository) *AuthService {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-256-bit-secret" // デフォルト値（開発用）
	}
	return &AuthService{
		authRepo:    authRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: time.Hour * 24,
	}
}

// LoginRequest ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest ユーザー登録リクエスト
type RegisterRequest struct {
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=6"`
	DisplayName *string `json:"display_name" validate:"omitempty,max=100"`
}

// AuthResponse 認証レスポンス
type AuthResponse struct {
	User struct {
		ID          string  `json:"id"`
		Email       string  `json:"email"`
		DisplayName *string `json:"display_name"`
		CreatedAt   string  `json:"created_at"`
	} `json:"user"`
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

// Register ユーザー登録
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	ctx := context.Background()

	// メール重複チェック
	existingUser, err := s.authRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// パスワードハッシュ化
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// ユーザー作成
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		DisplayName:  req.DisplayName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.authRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// JWT生成
	token, expiresIn, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	response := &AuthResponse{
		Token:     token,
		ExpiresIn: expiresIn,
	}
	response.User.ID = user.ID
	response.User.Email = user.Email
	response.User.DisplayName = user.DisplayName
	response.User.CreatedAt = user.CreatedAt.Format("2006-01-02T15:04:05Z")

	return response, nil
}

// Login ログイン
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	ctx := context.Background()

	// ユーザー取得
	user, err := s.authRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// パスワード検証
	if !crypto.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// JWT生成
	token, expiresIn, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	response := &AuthResponse{
		Token:     token,
		ExpiresIn: expiresIn,
	}
	response.User.ID = user.ID
	response.User.Email = user.Email
	response.User.DisplayName = user.DisplayName
	response.User.CreatedAt = user.CreatedAt.Format("2006-01-02T15:04:05Z")

	return response, nil
}

// generateToken JWT生成
func (s *AuthService) generateToken(user *models.User) (string, int64, error) {
	expiresAt := time.Now().Add(s.tokenExpiry)

	claims := &crypto.JWTClaims{
		Sub:   user.ID,
		Email: user.Email,
		Role:  "user", // デフォルトロール
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(s.tokenExpiry.Seconds()), nil
}
