CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('STUDENT', 'TEACHER', 'ADMIN'))
);

CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);

CREATE TABLE IF NOT EXISTS students (
    id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    fio VARCHAR(300) NOT NULL,
    group_name VARCHAR(30) NOT NULL,
    phone_number VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_students_group ON students(group_name);

CREATE TABLE IF NOT EXISTS teachers (
    id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    fio VARCHAR(300) NOT NULL,
    phone_number VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS teachers_group (
    teacher_id INT REFERENCES teachers(id) ON DELETE CASCADE,
    group_name VARCHAR(30) NOT NULL,
    PRIMARY KEY (teacher_id, group_name)
);