package user

import (
	"github.com/Artexus/api-matthew-backend/constant"
	dbUser "github.com/Artexus/api-matthew-backend/models/user"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type RepositoryInterface interface {
	Create(entity dbUser.User) (err error)
	ExistUsername(username string) (exist bool, err error)
	ExistEmail(email string) (exist bool, err error)
	TakeUserByIdentifier(identifier string) (entity dbUser.User, err error)
	TakeUserByIDs(ids []int) (entities []dbUser.User, err error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) Create(entity dbUser.User) (err error) {
	query := repo.db.Model(&entity).Begin().
		Create(&entity)
	if err = query.Error; err != nil {
		query.Rollback()
		return
	}

	err = query.Commit().Error
	return
}

func (repo *Repository) ExistUsername(username string) (exist bool, err error) {
	entity := dbUser.User{
		Username: username,
	}
	exist, err = repo.exist(entity, constant.IdentifierTypeUsername)
	errors.Wrap(err, "exist")
	return
}

func (repo *Repository) ExistEmail(email string) (exist bool, err error) {
	entity := dbUser.User{
		Email: email,
	}
	exist, err = repo.exist(entity, constant.IdentifierTypeEmail)
	errors.Wrap(err, "exist")
	return
}

func (repo *Repository) exist(entity dbUser.User, identifier string) (exist bool, err error) {
	count := 0
	query := repo.db.Model(&entity)

	if identifier == constant.IdentifierTypeEmail {
		query = query.Where("email = ?", entity.Email)
	} else if identifier == constant.IdentifierTypeUsername {
		query = query.Where("username = ? ", entity.Username)
	}

	query.Count(&count)
	if err = query.Error; err != nil {
		exist = false
		return
	}

	exist = count > 0
	return
}

func (repo *Repository) TakeUserByIdentifier(identifier string) (entity dbUser.User, err error) {
	query := repo.db.Model(&entity).
		Where("username = ?", identifier).
		Or("email = ?", identifier).
		Take(&entity)
	err = query.Error
	if query.RowsAffected == 0 {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (repo *Repository) TakeUserByIDs(ids []int) (entities []dbUser.User, err error) {
	query := repo.db.Model(&entities).
		Where("id IN (?)", ids).
		Find(&entities)
	err = query.Error
	return
}
