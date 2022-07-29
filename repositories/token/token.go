package token

import (
	dbToken "github.com/Artexus/api-matthew-backend/models/token"
	"github.com/jinzhu/gorm"
)

type RepositoryInterface interface {
	TakeByToken(token string, tokenType string) (entity dbToken.Token, err error)
	TakeByUser(userID int, tokenType string) (entity dbToken.Token, err error)
	ExistByUserID(userID int, tokenType string) (exist bool, err error)
	Create(entity dbToken.Token) (err error)
	Update(entity dbToken.Token) (err error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) ExistByUserID(userID int, tokenType string) (exist bool, err error) {
	count := 0
	err = repo.db.Model(&dbToken.Token{}).
		Where("user_id = ?", userID).
		Where("type = ?", tokenType).
		Count(&count).Error

	exist = count > 0
	return
}

func (repo *Repository) Create(entity dbToken.Token) (err error) {
	query := repo.db.Model(&entity).Begin().
		Create(&entity)
	if err = query.Error; err != nil {
		query.Rollback()
		return
	}

	err = query.Commit().Error
	return
}

func (repo *Repository) Update(entity dbToken.Token) (err error) {
	query := repo.db.Model(&entity).Begin().
		Where("user_id = ?", entity.UserID).
		Where("type = ?", entity.Type).
		Update(&entity)

	if err = query.Error; err != nil {
		query.Rollback()
		return
	}

	err = query.Commit().Error
	return
}

func (repo *Repository) TakeByToken(token string, tokenType string) (entity dbToken.Token, err error) {
	query := repo.db.Model(&entity).
		Where("token = ?", token).
		Where("type = ?", tokenType).
		Take(&entity)

	if query.RowsAffected < 0 {
		err = gorm.ErrRecordNotFound
		return
	}

	err = query.Error
	return
}

func (repo *Repository) TakeByUser(userID int, tokenType string) (entity dbToken.Token, err error) {
	query := repo.db.Model(&entity).
		Where("user_id = ?", userID).
		Where("type = ?", tokenType).
		Take(&entity)

	if query.RowsAffected < 0 {
		err = gorm.ErrRecordNotFound
		return
	}

	err = query.Error
	return
}
