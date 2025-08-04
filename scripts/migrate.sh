#!/bin/bash

# データベースマイグレーション実行スクリプト

set -e

# 環境変数の確認
if [ -z "$DB_HOST" ] || [ -z "$DB_NAME" ] || [ -z "$DB_USER" ]; then
    echo "エラー: 環境変数が設定されていません"
    echo "DB_HOST, DB_NAME, DB_USER, DB_PASSWORD を設定してください"
    exit 1
fi

# PostgreSQL接続情報
PGHOST=$DB_HOST
PGPORT=${DB_PORT:-5432}
PGDATABASE=$DB_NAME
PGUSER=$DB_USER
PGPASSWORD=$DB_PASSWORD

export PGHOST PGPORT PGDATABASE PGUSER PGPASSWORD

echo "🚀 データベースマイグレーション開始..."
echo "接続先: $PGHOST:$PGPORT/$PGDATABASE"

# マイグレーションファイルのディレクトリ
MIGRATION_DIR="infra/database/migrations"

# マイグレーションファイルを順番に実行
for file in $MIGRATION_DIR/*.sql; do
    if [ -f "$file" ]; then
        echo "📄 実行中: $(basename $file)"
        psql -f "$file" || {
            echo "❌ エラー: $(basename $file) の実行に失敗しました"
            exit 1
        }
    fi
done

echo "✅ マイグレーション完了！"

# テーブル一覧を表示
echo "📋 作成されたテーブル:"
psql -c "\dt"