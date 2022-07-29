package user

import (
	httpUser "github.com/Artexus/api-matthew-backend/http/user"
	"github.com/Artexus/api-matthew-backend/repositories/user"
	"github.com/pkg/errors"
)

type Service struct {
	repo user.RepositoryInterface
}

var instance *Service

func NewService(repo user.RepositoryInterface) *Service {
	if instance == nil {
		instance = &Service{
			repo: repo,
		}
	}
	return instance
}

func (svc *Service) GetUserByIDs(ids []int) (resp map[int]httpUser.UserResponse, err error) {
	resp = make(map[int]httpUser.UserResponse, 0)
	entities, err := svc.repo.TakeUserByIDs(ids)
	if err != nil {
		err = errors.Wrap(err, "db: take user by ids")
		return
	}

	for _, entity := range entities {
		resp[int(entity.ID)] = httpUser.UserResponse{
			EncID:    entity.EncID(),
			Username: entity.Username,
			Email:    entity.Email,
		}
	}
	return
}
