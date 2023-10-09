package user

type CreateUserDTO struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
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
