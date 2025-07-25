package login

type LoginCommand struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
