-- 企業マスタテーブル
CREATE TABLE companies (
  -- 基本情報
  id SERIAL PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  industry VARCHAR(100),
  size_category VARCHAR(50) CHECK (size_category IN ('startup', 'medium', 'large')),
  location VARCHAR(200),
  website_url TEXT,
  
  -- 個人メモ
  notes TEXT,
  rating INTEGER CHECK (rating BETWEEN 1 AND 5),
  
  -- メタデータ
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE -- 論理削除
);

-- インデックス
CREATE INDEX idx_companies_name ON companies(name);
CREATE INDEX idx_companies_industry ON companies(industry);
CREATE INDEX idx_companies_deleted ON companies(deleted_at);

-- updated_atトリガー
CREATE TRIGGER update_companies_updated_at 
  BEFORE UPDATE ON companies 
  FOR EACH ROW 
  EXECUTE FUNCTION update_updated_at_column();