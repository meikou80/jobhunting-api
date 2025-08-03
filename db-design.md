# シンプル転職活動管理システム DB設計書

## 基本情報
- **データベース**: PostgreSQL (Supabase)
- **テーブル数**: 4テーブル（勤怠管理と同等）
- **設計思想**: シンプル・確実性重視
- **認証**: Supabase Auth (Row Level Security有効)

## テーブル設計

### 1. companies（企業マスタ）

```sql
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
```

### 2. jobs（求人情報）

```sql
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
```

### 3. applications（応募記録）

```sql
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
```

### 4. user_profiles（ユーザー設定）

```sql
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
```

## データ統合仕様

### activity_log フィールド（JSON形式）
```json
{
  "events": [
    {
      "date": "2024-01-16",
      "type": "applied",
      "note": "応募書類提出"
    },
    {
      "date": "2024-01-25", 
      "type": "interview",
      "interviewer": "田中部長",
      "note": "技術面接。Goについて詳しく質問された",
      "result": "passed"
    }
  ]
}
```

### default_keywords フィールド
```sql
-- カンマ区切りのキーワード文字列
UPDATE user_profiles 
SET default_keywords = 'Go,Golang,バックエンド,リモート'
WHERE user_id = $1;
```

## Row Level Security (RLS)

```sql
-- 基本的なRLS設定
ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
ALTER TABLE jobs ENABLE ROW LEVEL SECURITY;
ALTER TABLE applications ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_profiles ENABLE ROW LEVEL SECURITY;

-- 企業・求人は全ユーザー閲覧可
CREATE POLICY "All users can view companies" ON companies FOR SELECT USING (true);
CREATE POLICY "All users can view jobs" ON jobs FOR SELECT USING (true);

-- 応募情報は自分のもののみ
CREATE POLICY "Users can manage own applications" ON applications USING (auth.uid() = user_id);
CREATE POLICY "Users can manage own profile" ON user_profiles USING (auth.uid() = user_id);
```

## 関数・トリガー（勤怠管理と同じパターン）

```sql
-- updated_at自動更新
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

-- 各テーブルにトリガー適用
CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_jobs_updated_at BEFORE UPDATE ON jobs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_applications_updated_at BEFORE UPDATE ON applications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## パフォーマンス設計

### インデックス戦略（勤怠管理と同様にシンプル）
- 検索用: `jobs.title` (全文検索)
- フィルタ用: `salary_min/max`, `location`
- 関連検索: `company_id`, `user_id`
- 日付範囲: `posted_date`, `applied_date`

### クエリ最適化
```sql
-- よく使用されるクエリのビュー
CREATE VIEW my_applications AS
SELECT 
  a.*,
  j.title,
  c.name as company_name,
  j.salary_min,
  j.salary_max
FROM applications a
JOIN jobs j ON a.job_id = j.id
JOIN companies c ON j.company_id = c.id
WHERE a.user_id = auth.uid()
AND a.deleted_at IS NULL;
```

この設計により、勤怠管理アプリと同等のシンプルさを保ちながら、転職活動の核心機能を効率的に管理できます。複雑な機能は後から必要に応じて追加することも可能です。