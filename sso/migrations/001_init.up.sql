CREATE TABLE roles (
  id        int GENERATED ALWAYS AS IDENTITY  PRIMARY KEY,   -- Автоинкрементируемый первичный ключ
  name      TEXT NOT NULL            -- Название роли
);

CREATE TABLE users (
    id        bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,   -- Автоинкрементируемый первичный ключ
    email     TEXT UNIQUE,    -- Уникальный email
    phone     TEXT UNIQUE,    -- Уникальный номер телефона
    pass_hash BYTEA NOT NULL,           -- Пароль в виде хэша (байты)
    role_id   int REFERENCES roles(id)
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);

CREATE TABLE apps (
  id        int GENERATED ALWAYS AS IDENTITY,   -- Автоинкрементируемый первичный ключ
  name      TEXT NOT NULL,           -- Название приложения
  secret    TEXT NOT NULL            -- Секретное слово приложения
);

INSERT INTO roles (name) VALUES ('admin');
INSERT INTO roles (name) VALUES ('user');
INSERT INTO apps (name, secret) VALUES ('web', 'web_secret');