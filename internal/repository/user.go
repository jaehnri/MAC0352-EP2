package repository

type Interruption interface {
	Create(name string, password string)
	ChangePassword(name string, password string)
	Login(name string)
	Logout(name string)
}
