CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('STUDENT', 'TEACHER', 'ADMIN'))


    CONSTRAINT chk_email_format CHECK (email ~ '^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$'),

    CONSTRAINT chk_password_format CHECK (length(password) >= 8 AND password ~ '[a-zA-Z]' AND password ~ '[0-9]')
);


CREATE TABLE IF NOT EXISTS email_verification (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    token UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_email_verifications_user_id ON email_verification(user_id);
CREATE INDEX IF NOT EXISTS idx_email_verifications_token ON email_verification(token);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS students (
    id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    fio VARCHAR(300) NOT NULL,
    group_name VARCHAR(30) NOT NULL,
    phone_number VARCHAR(50),

    CONSTRAINT chk_student_phone CHECK (
        phone_number IS NULL OR phone_number ~ '^\+?[1-9]\d{1,14}$'
    ),

    CONSTRAINT chk_student_fio CHECK (fio ~ '^[a-zA-Zа-яА-ЯёЁ\s\-]+$')
);

CREATE INDEX IF NOT EXISTS idx_students_group ON students(group_name);

CREATE TABLE IF NOT EXISTS teachers (
    id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    fio VARCHAR(300) NOT NULL,
    phone_number VARCHAR(50),

    CONSTRAINT chk_teacher_phone CHECK (
        phone_number IS NULL OR phone_number ~ '^\+?[1-9]\d{1,14}$'
    ),

    CONSTRAINT chk_teacher_fio CHECK (fio ~ '^[a-zA-Zа-яА-ЯёЁ\s\-]+$')
);

CREATE TABLE IF NOT EXISTS teachers_group (
    teacher_id INT REFERENCES teachers(id) ON DELETE CASCADE,
    group_name VARCHAR(30) NOT NULL,
    PRIMARY KEY (teacher_id, group_name)
);