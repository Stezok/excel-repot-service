package service

type UpdateTimeRepository interface {
	GetLastUpdateTime() (int64, error)
	SetLastUpdateTime(int64) error
}

type DefaultUpdateTimeService struct {
	repo UpdateTimeRepository
}

func (uts *DefaultUpdateTimeService) GetLastUpdateTime() (int64, error) {
	return uts.repo.GetLastUpdateTime()
}

func (uts *DefaultUpdateTimeService) SetLastUpdateTime(time int64) error {
	return uts.repo.SetLastUpdateTime(time)
}

func NewDefaultUpdateTimeService(repo UpdateTimeRepository) *DefaultUpdateTimeService {
	return &DefaultUpdateTimeService{
		repo: repo,
	}
}
