package auth

type AuthImpl struct{}

func NewAuthImpl() (*AuthImpl, error) {
	return &AuthImpl{}, nil
}
