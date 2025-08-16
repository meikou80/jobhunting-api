-- Row Level Security (RLS) 設定
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