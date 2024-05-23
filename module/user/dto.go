package user

import "time"

// UserDTO 表示输出到外部的用户信息

type UserDTO struct {
	ID          int64     `json:"id,omitempty"`
	Username    string    `json:"username,omitempty"`
	Nickname    string    `json:"nickname,omitempty"`
	Avatar      string    `json:"avatar,omitempty"`
	Description string    `json:"description,omitempty"`
	CreateAt    time.Time `json:"created_at"`
}
