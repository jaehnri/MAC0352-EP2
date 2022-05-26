package services

type User interface {
	Create(name string, password string)
	ChangePassword(name string, password string)
	Login(name string)
	Logout(name string)
}

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) Create(name string, password string) {
	// TODO: Implement user creation
}

func (u *UserService) ChangePassword(name string, password string) {
	// TODO: Implement password change
}

func (u *UserService) Login(name string) {
	// TODO: Implement login
}

func (u *UserService) Logout(name string) {
	// TODO: Implement logout
}
