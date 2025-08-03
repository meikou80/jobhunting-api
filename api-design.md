# シンプル転職活動管理システム API設計書

## 基本情報
- **プロトコル**: HTTP/HTTPS
- **データ形式**: JSON
- **認証**: Supabase Auth (JWT)
- **設計思想**: 勤怠管理アプリと同等のシンプルさ

## エンドポイント一覧

### 1. 企業管理API (/api/companies)

#### GET /api/companies
**概要**: 企業一覧取得

**クエリパラメータ**:
```
search: string (企業名検索)
industry: string (業界フィルタ)
size: string (企業規模)
page: integer (default: 1)
limit: integer (default: 20)
```

**レスポンス**:
```json
{
  "companies": [
    {
      "id": 1,
      "name": "株式会社テック",
      "industry": "IT・インターネット",
      "size_category": "startup",
      "location": "東京都渋谷区",
      "rating": 4
    }
  ],
  "total": 45,
  "page": 1
}
```

#### POST /api/companies
**概要**: 企業情報登録

**リクエスト**:
```json
{
  "name": "株式会社テック",
  "industry": "IT・インターネット",
  "size_category": "startup",
  "location": "東京都渋谷区",
  "website_url": "https://company.com",
  "notes": "興味のある会社",
  "rating": 4
}
```

### 2. 求人管理API (/api/jobs)

#### GET /api/jobs
**概要**: 求人一覧取得（検索機能付き）

**クエリパラメータ**:
```
keyword: string (タイトル検索)
company_id: integer (企業ID)
location: string (勤務地)
salary_min: integer (最低年収)
employment_type: string (雇用形態)
remote_option: string (リモート可否)
applied_status: string (not_applied, applied, all)
page: integer
limit: integer
```

**レスポンス**:
```json
{
  "jobs": [
    {
      "id": 1,
      "title": "Go言語エンジニア",
      "company": {
        "id": 1,
        "name": "株式会社テック"
      },
      "salary_range": "600万円〜1000万円",
      "location": "東京都/リモート可",
      "employment_type": "正社員",
      "posted_date": "2024-01-15",
      "application_status": "not_applied", // not_applied, applied
      "my_priority": null // 応募している場合の優先度
    }
  ],
  "total": 156,
  "page": 1
}
```

#### GET /api/jobs/:id
**概要**: 求人詳細取得

**レスポンス**:
```json
{
  "job": {
    "id": 1,
    "title": "Go言語エンジニア",
    "description": "詳細な求人説明...",
    "company": {
      "id": 1,
      "name": "株式会社テック",
      "industry": "IT・インターネット",
      "notes": "個人メモ"
    },
    "salary_min": 6000000,
    "salary_max": 10000000,
    "location": "東京都渋谷区",
    "employment_type": "正社員",
    "remote_option": "リモート可",
    "source_site": "wantedly",
    "source_url": "https://wantedly.com/projects/123",
    "posted_date": "2024-01-15",
    "deadline_date": "2024-02-15",
    "application": {
      "status": "applied",
      "applied_date": "2024-01-16",
      "priority": 2,
      "notes": "応募メモ"
    }
  }
}
```

#### POST /api/jobs
**概要**: 求人手動登録

**リクエスト**:
```json
{
  "company_id": 1,
  "title": "Go言語エンジニア",
  "description": "求人説明",
  "salary_min": 6000000,
  "salary_max": 10000000,
  "location": "東京都渋谷区",
  "employment_type": "正社員",
  "remote_option": "リモート可",
  "source_url": "https://example.com/job/123",
  "deadline_date": "2024-02-15"
}
```

### 3. 応募管理API (/api/applications)

#### GET /api/applications
**概要**: 応募一覧取得

**クエリパラメータ**:
```
status: string (応募ステータス)
priority: integer (優先度)
company_id: integer (企業ID)
sort: string (applied_date_desc, priority_asc, updated_desc)
```

**レスポンス**:
```json
{
  "applications": [
    {
      "id": 1,
      "job": {
        "id": 1,
        "title": "Go言語エンジニア",
        "company_name": "株式会社テック"
      },
      "status": "interview",
      "status_display": "面接中",
      "applied_date": "2024-01-16",
      "priority": 2,
      "priority_display": "中",
      "notes": "書類選考通過",
      "next_action": "2次面接が1/25に予定",
      "updated_at": "2024-01-20T15:30:00Z"
    }
  ],
  "summary": {
    "total": 15,
    "by_status": {
      "interested": 3,
      "applied": 5,
      "interview": 4,
      "offer": 1,
      "rejected": 2
    }
  }
}
```

#### POST /api/applications
**概要**: 応募登録

**リクエスト**:
```json
{
  "job_id": 1,
  "status": "applied",
  "applied_date": "2024-01-16",
  "priority": 2,
  "notes": "応募時のメモ"
}
```

#### PUT /api/applications/:id
**概要**: 応募状況更新

**リクエスト**:
```json
{
  "status": "interview",
  "priority": 1,
  "notes": "1次面接通過！2次面接は来週",
  "activity_log": {
    "events": [
      {
        "date": "2024-01-25",
        "type": "interview",
        "interviewer": "田中部長",
        "note": "技術面接。Goの経験について詳しく質問された",
        "result": "passed"
      }
    ]
  }
}
```

#### GET /api/applications/:id
**概要**: 応募詳細取得

**レスポンス**:
```json
{
  "application": {
    "id": 1,
    "job": {
      "id": 1,
      "title": "Go言語エンジニア",
      "company": {
        "id": 1,
        "name": "株式会社テック"
      }
    },
    "status": "interview",
    "applied_date": "2024-01-16",
    "priority": 1,
    "notes": "応募メモ",
    "interview_notes": "面接での印象など",
    "company_research": "企業研究メモ",
    "activity_log": {
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
          "note": "技術面接",
          "result": "passed"
        }
      ]
    }
  }
}
```

### 4. ユーザー設定API (/api/profile)

#### GET /api/profile
**概要**: ユーザー設定取得

**レスポンス**:
```json
{
  "profile": {
    "display_name": "田中太郎",
    "desired_salary_min": 6000000,
    "desired_salary_max": 10000000,
    "desired_location": "東京都・リモート",
    "desired_employment_type": "正社員",
    "default_keywords": "Go,Golang,バックエンド",
    "excluded_companies": "ブラック企業A,ブラック企業B",
    "email_notifications": true
  }
}
```

#### PUT /api/profile
**概要**: ユーザー設定更新

**リクエスト**:
```json
{
  "display_name": "田中太郎",
  "desired_salary_min": 6000000,
  "desired_salary_max": 10000000,
  "desired_location": "東京都・リモート",
  "default_keywords": "Go,Golang,バックエンド,クラウド",
  "email_notifications": true
}
```

### 5. 簡易スクレイピングAPI (/api/scraping)

#### POST /api/scraping/run
**概要**: 手動スクレイピング実行

**リクエスト**:
```json
{
  "keywords": ["Go", "Golang"],
  "location": "東京",
  "sites": ["wantedly", "green"]
}
```

**レスポンス**:
```json
{
  "message": "スクレイピングを開始しました",
  "estimated_time": "約3分",
  "job_id": "batch-20240116-001"
}
```

#### GET /api/scraping/status/:job_id
**概要**: スクレイピング状況確認

**レスポンス**:
```json
{
  "status": "completed",
  "progress": "100%",
  "results": {
    "found": 23,
    "new": 8,
    "updated": 2
  },
  "completed_at": "2024-01-16T10:05:00Z"
}
```

### 6. ダッシュボードAPI (/api/dashboard)

#### GET /api/dashboard
**概要**: ダッシュボード情報取得

**レスポンス**:
```json
{
  "summary": {
    "total_applications": 15,
    "active_interviews": 3,
    "offers": 1,
    "this_week_activities": 5
  },
  "recent_jobs": [
    {
      "id": 1,
      "title": "新着求人タイトル",
      "company_name": "株式会社テック",
      "posted_date": "2024-01-16"
    }
  ],
  "upcoming_actions": [
    {
      "type": "interview",
      "company": "株式会社テック",
      "date": "2024-01-25",
      "note": "2次面接"
    }
  ]
}
```

## データ形式仕様

### activity_log の構造
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
      "note": "技術面接",
      "result": "passed"
    }
  ]
}
```

### キーワード管理
```javascript
// フロントエンドでの検索条件管理例
const searchConditions = {
  keywords: ['Go', 'Golang'],
  location: '東京',
  salaryMin: 6000000
};
localStorage.setItem('searchConditions', JSON.stringify(searchConditions));
```

## エラーハンドリング

### 標準エラーレスポンス（勤怠管理と同様）
```json
{
  "error": "エラーメッセージ",
  "code": "VALIDATION_ERROR"
}
```

### HTTPステータスコード
- `200 OK`: 成功
- `201 Created`: 作成成功  
- `400 Bad Request`: リクエスト不正
- `401 Unauthorized`: 認証失敗
- `404 Not Found`: リソース未発見
- `500 Internal Server Error`: サーバーエラー

この設計により、勤怠管理アプリと同等のシンプルさを保ちながら、転職活動の核心機能を効率的に管理できます。