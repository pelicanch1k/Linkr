package dto

// Блокировка пользователя
type BlockUserRequest struct {
	Blocked bool `json:"blocked"`
}

// Изменение роли
type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// Статистика пользователя
type UserStats struct {
	LastLogin    string `json:"last_login" db:"last_login"`
	LoginCount   int    `json:"login_count" db:"login_count"`
	SessionCount int    `json:"session_count" db:"session_count"`
}

// Системная статистика
type SystemStats struct {
	UserCount       int `json:"user_count"`
	ActiveUsers     int `json:"active_users"`
	LockedAccounts  int `json:"locked_accounts"`
	DeletedAccounts int `json:"deleted_accounts"`
}
