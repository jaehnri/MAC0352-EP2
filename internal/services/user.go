package services

type UserData struct {
	Username string
	State    string
	Points   int
}

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

func (u *UserService) Create(name string, password string) error {
	// TODO: Implement user creation
	return nil
}

func (u *UserService) ChangePassword(name string, oldpassword, newpassword string) error {
	// TODO: Implement password change
	return nil
}

func (u *UserService) Login(name string, password string) error {
	// TODO: Implement login
	return nil
}

func (u *UserService) Logout(name string) {
	// TODO: Implement logout
}

func (u *UserService) ListConnected() []UserData {
	// TODO: Implement logout
	return nil
}

func (u *UserService) ListAll() []UserData {
	// TODO: Implement logout
	return nil
}
