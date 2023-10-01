package user

type CreateUserDTO struct {
	Name        string `json:"name"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	OneTimeCode string `json:"one_time_code"` // is needed?
}

type DeleteUserDTO struct {
	ID string `json:"id"`
}
