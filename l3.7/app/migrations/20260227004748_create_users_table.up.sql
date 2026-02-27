CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'manager', 'viewer')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
INSERT INTO users (id, username, role) VALUES 
    ('cd2c0e4c-90ad-4172-b9a4-f27ef6db560c','admin', 'admin'),
    ('582e758c-5a4c-4f10-8176-349478854a92', 'manager', 'manager'),
    ('ae275fa5-aa0b-4881-9fc1-2ba18308070b', 'viewer', 'viewer');