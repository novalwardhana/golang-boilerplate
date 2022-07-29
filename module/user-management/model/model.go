package model

type Result struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	ID       int    `json:"id" gorm:"id"`
	Name     string `json:"name" gorm:"name"`
	Email    string `json:"email" gorm:"email"`
	Password string `json:"password" gorm:"password"`
}

type Role struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type UserHasRole struct {
	UserID int `json:"user_id" gorm:"user_id"`
	RoleID int `json:"role_id" gorm:"role_id"`
}

type NewUser struct {
	User
	Roles []int `json:"roles"`
}

type UserWithRoles struct {
	ID        int    `gorm:"id" json:"id"`
	Name      string `gorm:"name" json:"name"`
	Email     string `gorm:"email" json:"email"`
	Roles     []byte `gorm:"roles" json:"-"`
	JsonRoles []Role `gorm:"-" json:"roles"`
}

type Pagination struct {
	Page         int             `json:"page"`
	Limit        int             `json:"limit"`
	TotalData    int             `json:"total_data"`
	NumberOfPage int             `json:"number_of_page"`
	Data         []UserWithRoles `json:"data"`
}
