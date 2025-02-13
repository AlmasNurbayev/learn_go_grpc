-- Удаляем индексы
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_apps_name;

-- Удаляем таблицы
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS apps;
DROP TABLE IF EXISTS roles;