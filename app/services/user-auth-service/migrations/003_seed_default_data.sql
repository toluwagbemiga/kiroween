-- Insert default permissions
INSERT INTO permissions (name, resource, action, description) VALUES
    ('users:read', 'users', 'read', 'View user information'),
    ('users:write', 'users', 'write', 'Create and update users'),
    ('users:delete', 'users', 'delete', 'Delete users'),
    ('roles:read', 'roles', 'read', 'View roles'),
    ('roles:write', 'roles', 'write', 'Create and update roles'),
    ('roles:delete', 'roles', 'delete', 'Delete roles'),
    ('permissions:read', 'permissions', 'read', 'View permissions'),
    ('billing:read', 'billing', 'read', 'View billing information'),
    ('billing:write', 'billing', 'write', 'Manage billing and subscriptions'),
    ('analytics:read', 'analytics', 'read', 'View analytics data'),
    ('features:read', 'features', 'read', 'View feature flags'),
    ('features:write', 'features', 'write', 'Manage feature flags')
ON CONFLICT (name) DO NOTHING;

-- Insert default roles
INSERT INTO roles (name, description, is_system) VALUES
    ('admin', 'System administrator with full access', true),
    ('member', 'Standard team member', true),
    ('viewer', 'Read-only access', true)
ON CONFLICT (name) DO NOTHING;

-- Assign permissions to admin role (all permissions)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

-- Assign permissions to member role (read/write for most resources)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'member'
AND p.name IN ('users:read', 'billing:read', 'analytics:read', 'features:read')
ON CONFLICT DO NOTHING;

-- Assign permissions to viewer role (read-only)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'viewer'
AND p.name IN ('users:read', 'roles:read', 'permissions:read', 'analytics:read', 'features:read')
ON CONFLICT DO NOTHING;
