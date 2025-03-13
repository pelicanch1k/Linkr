-- Таблица для хранения информации о пользователях
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY, -- Уникальный идентификатор пользователя

    username VARCHAR(50) UNIQUE NOT NULL, -- Уникальное имя пользователя
    email VARCHAR(100) UNIQUE NOT NULL, -- Уникальный email пользователя
    password_hash VARCHAR(255) NOT NULL, -- Хэш пароля пользователя

    first_name VARCHAR(50), -- Имя пользователя
    last_name VARCHAR(50), -- Фамилия пользователя
    bio TEXT, -- Краткая информация о пользователе
    profile_picture_url VARCHAR(255), -- Ссылка на аватарку

    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Дата и время регистрации
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Дата и время последнего обновления
);

-- Таблица для хранения постов пользователей
CREATE TABLE posts (
    post_id SERIAL PRIMARY KEY, -- Уникальный идентификатор поста
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, -- Связь с пользователем
    content TEXT NOT NULL, -- Текст поста
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Дата и время создания поста
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Дата и время последнего обновления
);

-- Таблица для хранения переписок между пользователями
CREATE TABLE messages (
    message_id SERIAL PRIMARY KEY, -- Уникальный идентификатор сообщения
    sender_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, -- Связь с отправителем
    receiver_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, -- Связь с получателем
    content TEXT NOT NULL, -- Текст сообщения
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Дата и время отправки сообщения
    is_read BOOLEAN DEFAULT FALSE -- Флаг, указывающий, прочитано ли сообщение
);

-- Таблица для хранения лайков под постами
CREATE TABLE likes (
    like_id SERIAL PRIMARY KEY, -- Уникальный идентификатор лайка
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, -- Связь с пользователем
    post_id INT NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE, -- Связь с постом
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Дата и время постановки лайка
);

-- Таблица для хранения комментариев под постами
CREATE TABLE comments (
    comment_id SERIAL PRIMARY KEY, -- Уникальный идентификатор комментария
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, -- Связь с пользователем
    post_id INT NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE, -- Связь с постом
    content TEXT NOT NULL, -- Текст комментария
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Дата и время создания комментария
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Дата и время последнего обновления
);