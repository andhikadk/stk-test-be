-- Seed default admin user
-- Created at: 2024-10-25
-- Password: admin123 (should be changed in production!)

-- Insert admin user if not exists
INSERT INTO users (name, email, password, role, is_active)
VALUES (
    'Admin User',
    'admin@example.com',
    '$2a$10$slYQmyNdGzin7olVN3VN2OPST9/PgBkqquzi.Ss8KIUgO2t0jWMUe', -- bcrypt hash of 'admin123'
    'admin',
    true
)
ON CONFLICT (email) DO NOTHING;
