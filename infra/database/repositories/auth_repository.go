package repositories

import (
	"context"
	"errors"

	"jobhunting-api/domain/models"
	"jobhunting-api/domain/repositories"

	"gorm.io/gorm"
)

// authRepository 認証リポジトリの実装
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository 認証リポジトリのコンストラクタ
func NewAuthRepository(db *gorm.DB) repositories.AuthRepository {
	return &authRepository{
		db: db,
	}
}

// Create ユーザー作成
func (r *authRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetByID IDでユーザー取得
func (r *authRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail メールアドレスでユーザー取得
func (r *authRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update ユーザー更新
func (r *authRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

// Delete ユーザー削除（論理削除）
func (r *authRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		return err
	}
	return nil
}
