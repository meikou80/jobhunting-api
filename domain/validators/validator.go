package validators

import (
	"github.com/go-playground/validator/v10"
)

// Validator バリデーター
var Validator *validator.Validate

// Init バリデーターを初期化
func Init() {
	Validator = validator.New()

	// カスタムバリデーションルールをここに追加
	// 例：独自の日本語形式チェック等
}

// ValidateStruct 構造体のバリデーション
func ValidateStruct(s interface{}) error {
	return Validator.Struct(s)
}

// ValidateVar 単一フィールドのバリデーション
func ValidateVar(field interface{}, tag string) error {
	return Validator.Var(field, tag)
}
