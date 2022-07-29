package post

import (
	"log"

	dbPost "github.com/Artexus/api-matthew-backend/models/post"
	"github.com/Artexus/api-matthew-backend/utils/pagination"
	"github.com/jinzhu/gorm"
)

type RepositoryInterface interface {
	Insert(entity dbPost.Post) (err error)
	TakePost(pgn pagination.Pagination) (entities []dbPost.Post, err error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) Insert(entity dbPost.Post) (err error) {
	query := repo.db.Model(&entity).Begin().
		Create(&entity)

	log.Println(query)
	if err = query.Error; err != nil {
		query.Rollback()
		return
	}

	err = query.Commit().Error
	return
}

func (repo *Repository) TakePost(pgn pagination.Pagination) (entities []dbPost.Post, err error) {
	query := repo.db.Model(&dbPost.Post{}).
		Offset(pgn.Offset).
		Limit(pgn.Limit).
		Order("created_at desc")

	err = query.Find(&entities).Error
	return
}
