package entities

type Account struct {
	Username string
	Password string
}

type User struct {
	Account    *Account
	Nickname   string
	ProfileUri string
}
