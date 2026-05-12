package album

type Service interface {
	GetAlbums() []Album
	GetAlbum(id int64) (Album, bool)
	CreateAlbum(a Album) Album
	UpdateAlbum(id int64, a Album) (Album, bool)
	DeleteAlbum(id int64) bool
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAlbums() []Album {
	return s.repo.GetAll()
}

func (s *service) GetAlbum(id int64) (Album, bool) {
	return s.repo.GetByID(id)
}

func (s *service) CreateAlbum(a Album) Album {
	return s.repo.Create(a)
}

func (s *service) UpdateAlbum(id int64, a Album) (Album, bool) {
	return s.repo.Update(id, a)
}

func (s *service) DeleteAlbum(id int64) bool {
	return s.repo.Delete(id)
}
