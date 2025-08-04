-- ユーザー設定テーブル
CREATE TABLE user_profiles (
  -- 基本情報
  user_id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  display_name VARCHAR(100),
  
  -- 転職希望条件（最小限）
  desired_salary_min INTEGER,
  desired_salary_max INTEGER,
  desired_location VARCHAR(200),
  desired_employment_type VARCHAR(50),
  
  -- 検索条件（job_searchesの代替）
  default_keywords TEXT, -- カンマ区切りのキーワード
  excluded_companies TEXT, -- 除外企業リスト
  
  -- 通知設定（最小限）
  email_notifications BOOLEAN DEFAULT true,
  
  -- メタデータ
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- updated_atトリガー
CREATE TRIGGER update_user_profiles_updated_at 
  BEFORE UPDATE ON user_profiles 
  FOR EACH ROW 
  EXECUTE FUNCTION update_updated_at_column();