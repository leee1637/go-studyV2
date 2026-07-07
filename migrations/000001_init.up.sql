CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('STUDENT', 'TEACHER', 'ADMIN'))


    CONSTRAINT chk_login_format CHECK (login ~ '^[a-zA-Z0-9]{3,}$'),

    CONSTRAINT chk_password_format CHECK (length(password) >= 8 AND password ~ '[a-zA-Z]' AND password ~ '[0-9]')
);

CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);

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