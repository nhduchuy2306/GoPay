package album

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() []Album
	GetByID(id int64) (Album, bool)
	Create(a Album) Album
	Update(id int64, a Album) (Album, bool)
	Delete(id int64) bool
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	if err := db.AutoMigrate(&Album{}); err != nil {
		panic(err)
	}
	return &repo{db: db}
}

func (r *repo) GetAll() []Album {
	var albums []Album
	if err := r.db.Find(&albums).Error; err != nil {
		return nil
	}
	return albums
}

func (r *repo) GetByID(id int64) (Album, bool) {
	var album Album
	if err := r.db.First(&album, id).Error; err != nil {
		return Album{}, false
	}
	return album, true
}

func (r *repo) Create(a Album) Album {
	if err := r.db.Create(&a).Error; err != nil {
		return Album{}
	}
	return a
}

func (r *repo) Update(id int64, a Album) (Album, bool) {
	var existing Album
	if err := r.db.First(&existing, id).Error; err != nil {
		return Album{}, false
	}

	existing.Title = a.Title
	existing.Artist = a.Artist
	existing.Price = a.Price

	if err := r.db.Save(&existing).Error; err != nil {
		return Album{}, false
	}
	return existing, true
}

func (r *repo) Delete(id int64) bool {
	var album Album
	if err := r.db.First(&album, id).Error; err != nil {
		return false
	}

	if err := r.db.Delete(&album).Error; err != nil {
		return false
	}
	return true
}
