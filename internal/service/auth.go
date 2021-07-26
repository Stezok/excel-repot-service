package service

type AuthRepository interface {
	Login(string) (bool, error)
}

type DefaultAuthService struct {
	repo AuthRepository
}

func (as *DefaultAuthService) Login(pass string) (bool, error) {
	return as.repo.Login(pass)
}

func NewDefaultAuthService(repo AuthRepository) *DefaultAuthService {
	return &DefaultAuthService{
		repo: repo,
	}
}
