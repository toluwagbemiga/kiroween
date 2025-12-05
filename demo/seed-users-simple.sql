-- Demo Users Seed Data
-- Password for all users: "Demo123!"
-- Bcrypt hash generated with cost 10

-- Insert demo users with bcrypt hashed passwords
-- Hash for "Demo123!": $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy

INSERT INTO users (email, password_hash, name, is_active, created_at, updated_at) VALUES
    ('admin@haunted.dev', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Admin User', true, NOW(), NOW()),
    ('member@haunted.dev', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Member User', true, NOW(), NOW()),
    ('john.doe@haunted.dev', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'John Doe', true, NOW(), NOW()),
    ('jane.smith@haunted.dev', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Jane Smith', true, NOW(), NOW()),
    ('viewer@haunted.dev', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Viewer User', true, NOW(), NOW()),
    ('guest@haunted.dev', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Guest User', true, NOW(), NOW())
ON CONFLICT (email) DO UPDATE 
SET password_hash = EXCLUDED.password_hash, 
    name = EXCLUDED.name, 
    updated_at = NOW();

-- Assign admin role
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.email = 'admin@haunted.dev' AND r.name = 'admin'
ON CONFLICT DO NOTHING;

-- Assign member roles
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.email IN ('member@haunted.dev', 'john.doe@haunted.dev', 'jane.smith@haunted.dev') 
AND r.name = 'member'
ON CONFLICT DO NOTHING;

-- Assign viewer roles
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.email IN ('viewer@haunted.dev', 'guest@haunted.dev') 
AND r.name = 'viewer'
ON CONFLICT DO NOTHING;

-- Display created users
SELECT 
    u.email,
    u.name,
    r.name as role,
    u.is_active
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.email LIKE '%@haunted.dev'
ORDER BY r.name, u.email;
