-- 求人統合管理テーブル（メインテーブル）
CREATE TABLE jobs (
  -- 基本情報
  id SERIAL PRIMARY KEY,
  company_name VARCHAR(200) NOT NULL,        -- 企業名（非正規化）
  position_title VARCHAR(300) NOT NULL,      -- 職種名
  description TEXT,                          -- 求人詳細
  requirements TEXT,                         -- 応募要件
  
  -- プラットフォーム情報（重要）
  source_platform VARCHAR(50) NOT NULL,     -- wantedly, green, rikunabi, direct
  external_id VARCHAR(200),                 -- プラットフォーム側のID
  source_url TEXT,                          -- 元URLリンク
  
  -- 勤務条件
  salary_min INTEGER,
  salary_max INTEGER,
  location VARCHAR(200),
  employment_type VARCHAR(50) CHECK (employment_type IN ('正社員', '契約社員', '業務委託')),
  remote_option VARCHAR(50) CHECK (remote_option IN ('リモート可', 'ハイブリッド', '出社必須')),
  
  -- 日程
  posted_date DATE,
  deadline_date DATE,
  is_active BOOLEAN DEFAULT true,
  
  -- 個人管理
  status VARCHAR(50) DEFAULT 'interested' CHECK (
    status IN ('interested', 'applied', 'interview', 'offer', 'rejected', 'withdrawn')
  ),
  priority INTEGER DEFAULT 3 CHECK (priority BETWEEN 1 AND 5),
  personal_notes TEXT,                       -- 個人的なメモ
  
  -- 重複管理
  duplicate_group_id UUID,                   -- 重複グループID
  is_primary BOOLEAN DEFAULT true,           -- 代表求人フラグ
  duplicate_confidence DECIMAL(3,2),         -- 重複確信度 (0.0-1.0)
  
  -- メタデータ
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE,
  
  -- 重複防止制約
  UNIQUE(source_platform, external_id)
);

-- インデックス
CREATE INDEX idx_jobs_company_name ON jobs(company_name);
CREATE INDEX idx_jobs_title ON jobs USING gin(to_tsvector('japanese', position_title));
CREATE INDEX idx_jobs_platform ON jobs(source_platform);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_salary ON jobs(salary_min, salary_max);
CREATE INDEX idx_jobs_location ON jobs(location);
CREATE INDEX idx_jobs_posted ON jobs(posted_date DESC);
CREATE INDEX idx_jobs_duplicate_group ON jobs(duplicate_group_id);
CREATE INDEX idx_jobs_active ON jobs(is_active, posted_date DESC);
CREATE INDEX idx_jobs_deleted ON jobs(deleted_at);

-- updated_atトリガー
CREATE TRIGGER update_jobs_updated_at 
  BEFORE UPDATE ON jobs 
  FOR EACH ROW 
  EXECUTE FUNCTION update_updated_at_column();