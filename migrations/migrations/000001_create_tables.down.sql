-- Сначала удаляем таблицы, которые зависят от других таблиц
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS posts;

-- Удаляем таблицы, связанные с пользователями
DROP TABLE IF EXISTS user_notifications;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS user_tokens;

-- В конце удаляем основную таблицу пользователей
DROP TABLE IF EXISTS users;
