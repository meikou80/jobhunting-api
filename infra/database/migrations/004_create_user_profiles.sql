-- ユーザー設定テーブル
CREATE TABLE user_profiles (
  -- 基本情報
  user_id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  display_name VARCHAR(100),
  
  -- 転職希望条件
  desired_salary_min INTEGER,
  desired_salary_max INTEGER,
  desired_location VARCHAR(200),
  desired_employment_type VARCHAR(50),
  
  -- 検索・フィルタ設定
  default_keywords TEXT,                    -- カンマ区切りキーワード
  excluded_companies TEXT,                  -- 除外企業リスト
  active_platforms TEXT,                    -- 利用中プラットフォーム（カンマ区切り）
  
  -- 機能設定
  duplicate_detection_enabled BOOLEAN DEFAULT true,
  email_notifications BOOLEAN DEFAULT true,
  auto_scraping_enabled BOOLEAN DEFAULT false,
  
  -- メタデータ
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- updated_atトリガー
CREATE TRIGGER update_user_profiles_updated_at 
  BEFORE UPDATE ON user_profiles 
  FOR EACH ROW 
  EXECUTE FUNCTION update_updated_at_column();