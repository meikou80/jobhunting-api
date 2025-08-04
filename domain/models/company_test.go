package models

import (
	"testing"
	"time"
)

func TestCompanyToResponse(t *testing.T) {
	// テスト用の企業データ
	company := &Company{
		ID:           1,
		Name:         "株式会社テスト",
		Industry:     strPtr("IT・インターネット"),
		SizeCategory: strPtr("startup"),
		Location:     strPtr("東京都渋谷区"),
		WebsiteURL:   strPtr("https://example.com"),
		Notes:        strPtr("テスト企業"),
		Rating:       intPtr(4),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Response変換テスト
	resp := company.ToResponse()

	// 検証
	if resp.ID != company.ID {
		t.Errorf("Expected ID %d, got %d", company.ID, resp.ID)
	}

	if resp.Name != company.Name {
		t.Errorf("Expected Name %s, got %s", company.Name, resp.Name)
	}

	if resp.Industry == nil || *resp.Industry != *company.Industry {
		t.Errorf("Expected Industry %s, got %v", *company.Industry, resp.Industry)
	}
}

func TestCompanyFilter(t *testing.T) {
	filter := CompanyFilter{
		Search:   strPtr("テスト"),
		Industry: strPtr("IT"),
		Size:     strPtr("startup"),
		Page:     1,
		Limit:    20,
	}

	// 基本的な値確認
	if filter.Page != 1 {
		t.Errorf("Expected Page 1, got %d", filter.Page)
	}

	if filter.Limit != 20 {
		t.Errorf("Expected Limit 20, got %d", filter.Limit)
	}
}

// テスト用ヘルパー関数
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
