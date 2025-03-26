package dto

import "mime/multipart"

// Профиль пользователя
type UserProfile struct {
    ID        int    `json:"id" db:"id"`
    Username  string `json:"username" db:"username"`
    Email     string `json:"email" db:"email"`
    FirstName string `json:"first_name" db:"first_name"`
    LastName  string `json:"last_name" db:"last_name"`
    Avatar    string `json:"avatar" db:"profile_picture_url"`
    Role      string `json:"role" db:"role"`
    CreatedAt string `json:"created_at" db:"created_at"`
}

// Обновление профиля
type UpdateProfileRequest struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Username  string `json:"username"`
}

// Обновление email
type UpdateEmailRequest struct {
    Email string `json:"email" binding:"required"`
}

// Загрузка аватара
type AvatarUpload struct {
    File *multipart.FileHeader `form:"avatar"`
}

// Уведомления
type Notification struct {
    ID        int    `json:"id" db:"id"`
    UserID    int    `json:"user_id" db:"user_id"`
    Message   string `json:"message" db:"message"`
    Read      bool   `json:"read" db:"read"`
    CreatedAt string `json:"created_at" db:"created_at"`
} 