package user

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
	OneTimeCode  string `json:"one_time_code"`
}
