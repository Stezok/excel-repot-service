package service

type UpdateTimeRepository interface {
	GetLastUpdateTime() (int64, error)
	SetLastUpdateTime(int64) error
}

type CashedUpdateTimeService struct {
	repo UpdateTimeRepository
}

func (uts *CashedUpdateTimeService) GetLastUpdateTime() (int64, error) {
	return uts.repo.GetLastUpdateTime()
}

func (uts *CashedUpdateTimeService) SetLastUpdateTime(time int64) error {
	return uts.repo.SetLastUpdateTime(time)
}

func NewCashedUpdateTimeService(repo UpdateTimeRepository) *CashedUpdateTimeService {
	return &CashedUpdateTimeService{
		repo: repo,
	}
}
