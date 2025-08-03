# 転職活動管理システム API

シンプルで効率的な転職活動管理のためのバックエンドAPI

## 🎯 プロジェクト概要

転職活動中の求職者が直面する課題を解決するWebアプリケーションのバックエンドです。

### 解決する課題
- 情報散在: 複数の求人サイトに情報が分散
- 応募状況の混乱: 進捗状況が不明確
- 企業情報の蓄積不足: 研究内容が記録されない
- 進捗の見える化不足: 全体状況が把握できない

## 🛠️ 技術スタック

- **言語**: Go 1.21
- **フレームワーク**: Gin
- **データベース**: PostgreSQL
- **認証**: Supabase Auth
- **コンテナ**: Docker & Docker Compose

## 🚀 開発環境のセットアップ

### 前提条件
- Docker & Docker Compose
- Go 1.21+ (ローカル開発用)

### 1. プロジェクトのクローン
```bash
git clone <repository-url>
cd jobhunting-api
```

### 2. 環境変数の設定
```bash
cp .env.example .env
# .envファイルを編集して必要な値を設定
```

### 3. Dockerでの起動
```bash
# データベースとアプリケーションを起動
docker-compose up -d

# ログの確認
docker-compose logs -f app
```

### 4. ローカル開発（ホットリロード）
```bash
# Airのインストール（初回のみ）
go install github.com/cosmtrek/air@latest

# 依存関係のインストール
go mod tidy

# ホットリロードで起動
air
```

## 📚 API仕様

詳細なAPI仕様については `api-design.md` を参照してください。

### エンドポイント概要
- `/health` - ヘルスチェック
- `/api/v1/ping` - 疎通確認
- (今後追加予定)

## 🏗️ 開発ステータス

### Phase 1: 基盤構築 ✅
- [x] プロジェクト初期化
- [x] Docker環境構築
- [x] 基本API構造
- [ ] データベース設計
- [ ] 認証システム

### Phase 2: コア機能開発 (予定)
- [ ] 企業管理機能
- [ ] 求人管理機能
- [ ] 応募管理機能

## 🧪 テスト実行

```bash
# 単体テスト
go test ./...

# カバレッジ付きテスト
go test -cover ./...
```

## 📖 プロジェクト構成

```
jobhunting-api/
├── cmd/                    # アプリケーションエントリーポイント
├── domain/                 # ドメイン層
├── usecase/               # ユースケース層
├── interface/             # インターフェース層
├── infra/                 # インフラストラクチャ層
└── docs/                  # ドキュメント
```

## 🤝 コントリビューション

個人開発プロジェクトですが、フィードバックやアイデアをお待ちしています。

## 📝 ライセンス

MIT License