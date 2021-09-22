package service

type UpdateTimeRepository interface {
	GetLastUpdateTime(string) (int64, error)
	SetLastUpdateTime(string, int64) error
}

type DefaultUpdateTimeService struct {
	repo UpdateTimeRepository
}

func (uts *DefaultUpdateTimeService) GetLastUpdateTime(tag string) (int64, error) {
	return uts.repo.GetLastUpdateTime(tag)
}

func (uts *DefaultUpdateTimeService) SetLastUpdateTime(tag string, time int64) error {
	return uts.repo.SetLastUpdateTime(tag, time)
}

func NewDefaultUpdateTimeService(repo UpdateTimeRepository) *DefaultUpdateTimeService {
	return &DefaultUpdateTimeService{
		repo: repo,
	}
}
