-- Создаем расширение для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создаем таблицу пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу для refresh токенов
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid UUID NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    client_ip VARCHAR(45) NOT NULL,
    token_id VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_used BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_uuid);
CREATE INDEX idx_refresh_tokens_token_id ON refresh_tokens(token_id);

-- Добавляем тестовых пользователей
INSERT INTO users (name, username, password, email) VALUES
('Test User', 'testuser', 'password123', 'test@example.com'),
('John Doe', 'johndoe', 'password456', 'john@example.com'),
('Jane Smith', 'janesmith', 'password789', 'jane@example.com');

-- Получить UUID пользователей для тестирования
SELECT uuid FROM users WHERE username = 'testuser';