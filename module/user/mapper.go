package user

func ConvertUserToDTO(user *User) *UserDTO {
	if user == nil {
		return nil
	}
	return &UserDTO{
		ID:          user.Id,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Nickname,
		Description: user.Description,
		CreateAt:    user.CreateAt,
	}
}
