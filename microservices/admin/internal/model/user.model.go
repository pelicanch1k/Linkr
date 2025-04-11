package model

import "time"

// User представляет модель пользователя в системе.
type User struct {
    UserID          int       `db:"user_id"`           // Уникальный идентификатор пользователя
    Username        string    `db:"username"`           // Уникальное имя пользователя
    Email           string    `db:"email"`              // Уникальный email пользователя
    PasswordHash    string    `db:"password_hash"`      // Хэш пароля пользователя
    FirstName       string    `db:"first_name"`         // Имя пользователя
    LastName        string    `db:"last_name"`          // Фамилия пользователя
    Bio             *string    `db:"bio"`                // Краткая информация о пользователе
    ProfilePictureURL string  `db:"profile_picture_url"` // Ссылка на аватарку
    Role            string    `db:"role"`               // Роль пользователя (например, 'user', 'admin')
    Blocked         bool      `db:"blocked"`            // Статус заблокированности пользователя
    Deleted         bool      `db:"deleted"`            // Статус удаления пользователя
    DeletedAt       *time.Time `db:"deleted_at"`        // Время удаления пользователя (если удалён)
    LastLogin       *time.Time `db:"last_login"`        // Время последнего входа пользователя
    LoginCount      int       `db:"login_count"`        // Количество входов пользователя
    CreatedAt       time.Time `db:"created_at"`         // Дата и время регистрации
    UpdatedAt       time.Time `db:"updated_at"`         // Дата и время последнего обновления
}

