CREATE TYPE user_role AS ENUM ('user', 'admin', 'staff');

ALTER TABLE users 
ADD COLUMN first_name VARCHAR(100),
ADD COLUMN last_name VARCHAR(100),
ADD COLUMN phone VARCHAR(20),
ADD COLUMN role user_role DEFAULT 'user',
ADD COLUMN last_login_at TIMESTAMPTZ;

ALTER TABLE users 
ALTER COLUMN updated_at DROP DEFAULT;

ALTER TABLE users 
RENAME COLUMN password_hash TO password;

-- update old data 
UPDATE users
SET 
    first_name = 'N/A',
    last_name = 'N/A',
    phone = '0000000000',
    role = 'user'
WHERE 
    first_name IS NULL OR 
    last_name IS NULL OR 
    phone IS NULL OR 
    role IS NULL;

-- add not null
ALTER TABLE users
ALTER COLUMN first_name SET NOT NULL,
ALTER COLUMN last_name SET NOT NULL,
ALTER COLUMN phone SET NOT NULL,
ALTER COLUMN role SET NOT NULL;