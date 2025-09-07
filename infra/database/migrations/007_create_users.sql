-- 認証用ユーザーテーブル作成
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- インデックス作成
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- コメント追加
COMMENT ON TABLE users IS '認証用ユーザーテーブル';
COMMENT ON COLUMN users.id IS 'ユーザーID（UUID）';
COMMENT ON COLUMN users.email IS 'メールアドレス（ログイン用）';
COMMENT ON COLUMN users.password_hash IS 'ハッシュ化されたパスワード';
COMMENT ON COLUMN users.display_name IS '表示名';
COMMENT ON COLUMN users.created_at IS '作成日時';
COMMENT ON COLUMN users.updated_at IS '更新日時';
COMMENT ON COLUMN users.deleted_at IS '削除日時（論理削除）';
