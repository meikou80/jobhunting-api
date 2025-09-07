# 転職活動管理システム API

複数転職サービスの一元管理と重複応募防止を実現するバックエンドAPI

## 使用技術一覧

<p style="display: inline">
  <img src="https://img.shields.io/badge/-Go-00ADD8.svg?logo=go&style=for-the-badge&logoColor=white">
  <img src="https://img.shields.io/badge/-Gin-00ADD8.svg?logo=go&style=for-the-badge&logoColor=white">
  <img src="https://img.shields.io/badge/-PostgreSQL-336791.svg?logo=postgresql&style=for-the-badge&logoColor=white">
  <img src="https://img.shields.io/badge/-Docker-2496ED.svg?logo=docker&style=for-the-badge&logoColor=white">
  <img src="https://img.shields.io/badge/-JWT-000000.svg?logo=jsonwebtokens&style=for-the-badge&logoColor=white">
</p>

## 目次

1. [プロジェクトについて](#プロジェクトについて)
2. [技術選定理由](#技術選定理由)
3. [環境](#環境)
4. [セットアップ](#セットアップ)
5. [API動作確認](#api動作確認)

## プロジェクトについて

転職活動中の求職者が直面する課題を解決するWebアプリケーションのバックエンドAPI

### 解決する課題
- **重複応募リスク**: 同じ企業に異なるプラットフォームから重複応募してしまう
- **情報散在**: 複数の求人サイト（Wantedly、Green、doda等）に情報が分散
- **統合管理の不在**: 一元的な転職活動管理システムが存在しない

### 核心機能
- 複数プラットフォームからの求人一元管理
- 高精度な重複検知アルゴリズム（会社名・職種名・プラットフォーム情報）
- JWT認証によるセキュアなユーザー管理
- プラットフォーム別統計・ダッシュボード

## 技術選定理由

### バックエンド
GoのWebフレームワークはGinとEchoで迷ったが、フレームワークの中で一番スター数の多いGinを選択した。
Clean Architectureを採用することで保守性と拡張性を重視した設計とした。

### データベース
PostgreSQLを選択し、GORMを使用することでGo言語との親和性と型安全性を確保した。

### 認証
JWT + bcryptによるセキュアな認証システムを構築。
パスワードを平文のままDBに保存するのはセキュリティ的に良くないと考え、暗号化ライブラリを使用してハッシュ化して保存している。

### 気づいたこと/工夫したこと
複数の転職サービスからの重複応募を防ぐため、会社名・職種名・プラットフォーム情報を組み合わせた独自の重複検知アルゴリズムを実装した。
信頼度計算により、類似度の高い求人を効率的に検出できるシステムを構築した。

フロントエンドとバックエンドを完全に分離することでAPIを新しいフレームワークや技術にリプレイス（学習）したい時に変更を容易にできるようにした。

## 環境

| 言語・フレームワーク・ライブラリ | バージョン |
| ------------------------------ | ---------- |
| Go                             | 1.21       |
| Gin                            | 1.10.0     |
| PostgreSQL                     | 15         |
| GORM                           | 1.25.5     |

その他のパッケージのバージョンは go.mod を参照してください

## セットアップ

### 前提条件
- Docker & Docker Compose
- Go 1.21+

### 1. プロジェクトのクローン
```bash
git clone <repository-url>
cd jobhunting-api
```

### 2. Docker環境でのサーバー起動
```bash
# データベース起動
docker-compose up -d db

# 環境変数設定（.envファイルを作成し適切な値を設定）
cp .env.example .env

# サーバー起動
go run ./cmd
```

## API仕様

### 主要エンドポイント
- `POST /api/v1/auth/register` - ユーザー登録
- `POST /api/v1/auth/login` - ログイン
- `POST /api/v1/jobs` - 求人登録（重複チェック付き）
- `GET /api/v1/jobs` - 求人一覧取得（検索・フィルタ対応）
- `POST /api/v1/jobs/check-duplicates` - 重複チェック
- `GET /api/v1/jobs/dashboard` - ダッシュボード統計

### API仕様書
ブラウザで `http://localhost:8080/swagger/index.html` にアクセス

## 開発ステータス

### Phase 1: 基盤構築 ✅ 完了
- ✅ プロジェクト初期化・Docker環境構築
- ✅ データベース設計（PostgreSQL + マイグレーション）
- ✅ 認証システム（JWT + bcrypt）

### Phase 2: 求人統合管理 ✅ 完了
- ✅ 求人統合管理機能（プラットフォーム別登録・一覧・検索・フィルタ）
- ✅ 重複検知システム（会社名 + タイトル + 信頼度計算）
- ✅ プラットフォーム別統計・ダッシュボード
- ✅ API仕様書（Swagger UI）

**複数転職サービスの一元管理と重複防止** → **実装完了**