CREATE TABLE IF NOT EXISTS admins (
    id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    fio VARCHAR(300) NOT NULL,
    phone_number VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id INT PRIMARY KEY,
    id_users INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

CREATE TABLE IF NOT EXISTS registration_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fio VARCHAR(300) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('STUDENT', 'TEACHER')),
    phone_number VARCHAR(50),
    email VARCHAR(255) UNIQUE NOT NULL,
    group_name VARCHAR(30),
    token UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'expired')),
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_reg_requests_email ON registration_requests(email);
CREATE INDEX IF NOT EXISTS idx_reg_requests_token ON registration_requests(token);
CREATE INDEX IF NOT EXISTS idx_reg_requests_status ON registration_requests(status);