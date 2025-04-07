-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- UUID как первичный ключ
    name VARCHAR(50) NOT NULL,                     -- Имя пользователя
    age SMALLINT CHECK (age >= 0 AND age <= 120), -- Возраст (ограничение диапазона)
    gender VARCHAR(10) NOT NULL,                  -- Пол (male, female, other)
    email VARCHAR(255) UNIQUE NOT NULL,           -- Email (уникальный)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Дата создания
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  -- Дата обновления
);

CREATE INDEX idx_users_email ON users(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
