-- 求人情報テーブル
CREATE TABLE jobs (
  -- 基本情報
  id SERIAL PRIMARY KEY,
  company_id INTEGER REFERENCES companies(id),
  title VARCHAR(300) NOT NULL,
  description TEXT,
  
  -- 勤務条件
  salary_min INTEGER,
  salary_max INTEGER,
  location VARCHAR(200),
  employment_type VARCHAR(50) CHECK (employment_type IN ('正社員', '契約社員', '業務委託')),
  remote_option VARCHAR(50) CHECK (remote_option IN ('リモート可', 'ハイブリッド', '出社必須')),
  
  -- 外部情報
  source_site VARCHAR(100),
  source_url TEXT,
  external_id VARCHAR(200),
  
  -- 日程
  posted_date DATE,
  deadline_date DATE,
  is_active BOOLEAN DEFAULT true,
  
  -- メタデータ
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE,
  
  -- 重複防止
  UNIQUE(source_site, external_id)
);

-- インデックス
CREATE INDEX idx_jobs_company ON jobs(company_id);
CREATE INDEX idx_jobs_title ON jobs USING gin(to_tsvector('japanese', title));
CREATE INDEX idx_jobs_salary ON jobs(salary_min, salary_max);
CREATE INDEX idx_jobs_location ON jobs(location);
CREATE INDEX idx_jobs_posted ON jobs(posted_date DESC);
CREATE INDEX idx_jobs_active ON jobs(is_active, posted_date DESC);
CREATE INDEX idx_jobs_deleted ON jobs(deleted_at);

-- updated_atトリガー
CREATE TRIGGER update_jobs_updated_at 
  BEFORE UPDATE ON jobs 
  FOR EACH ROW 
  EXECUTE FUNCTION update_updated_at_column();