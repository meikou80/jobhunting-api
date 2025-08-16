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