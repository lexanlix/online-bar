package user

type CreateUserDTO struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	OneTimeCode string `json:"one_time_code"`
}

type PartUpdateUserDTO struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SignInUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DeleteUserDTO struct {
	ID string `json:"id"`
}

type RefreshUserDTO struct {
	RefreshToken string `json:"refresh_token"`
}
