-- 応募記録テーブル
CREATE TABLE applications (
  -- 基本情報
  id SERIAL PRIMARY KEY,
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
  job_id INTEGER REFERENCES jobs(id) ON DELETE CASCADE,
  company_id INTEGER REFERENCES companies(id),
  
  -- 応募情報
  status VARCHAR(50) NOT NULL DEFAULT 'interested' CHECK (
    status IN (
      'interested',    -- 興味あり
      'applied',       -- 応募済み
      'screening',     -- 書類選考中
      'interview',     -- 面接中
      'final',         -- 最終面接
      'offer',         -- 内定
      'rejected',      -- 不採用
      'withdrawn'      -- 辞退
    )
  ),
  applied_date DATE,
  priority INTEGER DEFAULT 3 CHECK (priority BETWEEN 1 AND 5), -- 1:高 5:低
  
  -- 活動記録（application_eventsの内容を統合）
  activity_log TEXT, -- JSON形式で面接履歴等を記録
  notes TEXT,
  interview_notes TEXT,
  company_research TEXT,
  
  -- 結果
  rejection_reason TEXT,
  offer_details TEXT, -- 内定条件をテキストで記録
  
  -- メタデータ
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE,
  
  -- 制約
  UNIQUE(user_id, job_id) -- 同じ求人への重複応募防止
);

-- インデックス
CREATE INDEX idx_applications_user ON applications(user_id);
CREATE INDEX idx_applications_user_status ON applications(user_id, status);
CREATE INDEX idx_applications_company ON applications(company_id);
CREATE INDEX idx_applications_applied ON applications(applied_date DESC);
CREATE INDEX idx_applications_priority ON applications(user_id, priority, updated_at DESC);
CREATE INDEX idx_applications_deleted ON applications(deleted_at);

-- updated_atトリガー
CREATE TRIGGER update_applications_updated_at 
  BEFORE UPDATE ON applications 
  FOR EACH ROW 
  EXECUTE FUNCTION update_updated_at_column();