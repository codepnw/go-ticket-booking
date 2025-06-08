ALTER TABLE users
DROP COLUMN first_name,
DROP COLUMN last_name,
DROP COLUMN phone,
DROP COLUMN role,
DROP COLUMN last_login_at;

ALTER TABLE users
ALTER COLUMN updated_at SET DEFAULT now();

ALTER TABLE users
RENAME COLUMN password TO password_hash;