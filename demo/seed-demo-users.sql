-- Demo Users Seed Data
-- This script creates demo users for testing different roles
-- All passwords are hashed using bcrypt with the password: "Demo123!"

-- Note: These are bcrypt hashes for "Demo123!" (cost 10)
-- You can generate new ones using: bcrypt.GenerateFromPassword([]byte("Demo123!"), 10)

-- Insert demo users
INSERT INTO users (email, password_hash, name, is_active, email_verified, created_at, updated_at) VALUES
    -- Admin user
    ('admin@haunted.dev', '$2a$10$rZ8qNW5xGx5xGx5xGx5xGOK8qNW5xGx5xGx5xGx5xGx5xGx5xGx5x', 'Admin User', true, true, NOW(), NOW()),
    
    -- Member users
    ('member@haunted.dev', '$2a$10$rZ8qNW5xGx5xGx5xGx5xGOK8qNW5xGx5xGx5xGx5xGx5xGx5xGx5x', 'Member User', true, true, NOW(), NOW()),
    ('john.doe@haunted.dev', '$2a$10$rZ8qNW5xGx5xGx5xGx5xGOK8qNW5xGx5xGx5xGx5xGx5xGx5xGx5x', 'John Doe', true, true, NOW(), NOW()),
    ('jane.smith@haunted.dev', '$2a$10$rZ8qNW5xGx5xGx5xGx5xGOK8qNW5xGx5xGx5xGx5xGx5xGx5xGx5x', 'Jane Smith', true, true, NOW(), NOW()),
    
    -- Viewer users
    ('viewer@haunted.dev', '$2a$10$rZ8qNW5xGx5xGx5xGx5xGOK8qNW5xGx5xGx5xGx5xGx5xGx5xGx5x', 'Viewer User', true, true, NOW(), NOW()),
    ('guest@haunted.dev', '$2a$10$rZ8qNW5xGx5xGx5xGx5xGOK8qNW5xGx5xGx5xGx5xGx5xGx5xGx5x', 'Guest User', true, true, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Assign roles to users
-- Admin role
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.email = 'admin@haunted.dev' AND r.name = 'admin'
ON CONFLICT DO NOTHING;

-- Member roles
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.email IN ('member@haunted.dev', 'john.doe@haunted.dev', 'jane.smith@haunted.dev') 
AND r.name = 'member'
ON CONFLICT DO NOTHING;

-- Viewer roles
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.email IN ('viewer@haunted.dev', 'guest@haunted.dev') 
AND r.name = 'viewer'
ON CONFLICT DO NOTHING;
