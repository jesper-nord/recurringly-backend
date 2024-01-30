package auth

type UserId uint

type Service interface {
	Login(username, password string) (*User, TokenPair, error)
	RegisterUser(username, password string) (*User, TokenPair, error)
	RefreshAccessToken(refreshToken string) (TokenPair, error)
}

type Repository interface {
	GetUser(userId UserId) (*User, error)
	FindUser(username string) (*User, error)
	SaveUser(user *User) (*User, error)
	Migrate() error
}
