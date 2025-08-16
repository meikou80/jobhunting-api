-- 応募履歴管理テーブル
CREATE TABLE applications (
  -- 基本情報
  id SERIAL PRIMARY KEY,
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
  job_id INTEGER REFERENCES jobs(id) ON DELETE CASCADE,
  
  -- 応募経路情報（重要）
  applied_via VARCHAR(50) NOT NULL,         -- 実際の応募プラットフォーム
  applied_date DATE NOT NULL,
  
  -- 応募状況
  status VARCHAR(50) NOT NULL DEFAULT 'applied' CHECK (
    status IN (
      'applied',        -- 応募済み
      'screening',      -- 書類選考中
      'interview',      -- 面接中
      'final',          -- 最終面接
      'offer',          -- 内定
      'rejected',       -- 不採用
      'withdrawn'       -- 辞退
    )
  ),
  current_stage VARCHAR(100),               -- 現在のステージ（1次面接、最終面接等）
  
  -- 活動記録
  activity_log TEXT,                        -- JSON形式の活動履歴
  notes TEXT,                               -- 応募メモ
  next_action TEXT,                         -- 次のアクション
  next_action_date DATE,                    -- 次のアクション予定日
  
  -- 結果
  rejection_reason TEXT,
  offer_details TEXT,                       -- 内定条件
  
  -- メタデータ
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE,
  
  -- 制約（同じ求人への重複応募防止）
  UNIQUE(user_id, job_id)
);

-- インデックス
CREATE INDEX idx_applications_user ON applications(user_id);
CREATE INDEX idx_applications_user_status ON applications(user_id, status);
CREATE INDEX idx_applications_applied_via ON applications(applied_via);
CREATE INDEX idx_applications_applied ON applications(applied_date DESC);
CREATE INDEX idx_applications_next_action ON applications(user_id, next_action_date);
CREATE INDEX idx_applications_deleted ON applications(deleted_at);

-- updated_atトリガー
CREATE TRIGGER update_applications_updated_at 
  BEFORE UPDATE ON applications 
  FOR EACH ROW 
  EXECUTE FUNCTION update_updated_at_column();