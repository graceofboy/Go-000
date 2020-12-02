package service

type Service struct {
	db *dao.DBService
}

func NewService() *Service {
	return &Service{db: dao.NewDBService()}
}

func (this *Service) FindPersion() ([]model.Persion, error) {
	return this.db.FindPersion()
}
