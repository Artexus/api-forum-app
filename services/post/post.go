package post

import (
	"github.com/Artexus/api-matthew-backend/constant"
	httpPost "github.com/Artexus/api-matthew-backend/http/post"
	dbPost "github.com/Artexus/api-matthew-backend/models/post"
	"github.com/Artexus/api-matthew-backend/repositories/post"
	"github.com/Artexus/api-matthew-backend/services/user"
	"github.com/Artexus/api-matthew-backend/utils/aes"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

type Service struct {
	repo    post.RepositoryInterface
	userSvc user.Service
}

var instance *Service

func NewService(repo post.RepositoryInterface, userSvc user.Service) *Service {
	if instance == nil {
		instance = &Service{
			repo:    repo,
			userSvc: userSvc,
		}
	}
	return instance
}

func (svc *Service) InsertPost(req httpPost.InsertPostRequest) (err error) {
	if len(req.Description) == 0 || len(req.Title) == 0 {
		err = constant.ErrDescriptionTitleEmpty
		return
	}

	if len(req.Description) > constant.MaxPostDescriptionLength {
		err = constant.ErrMaximumDescription
		return
	}

	entity := dbPost.Post{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
	}
	if req.FileID != nil {
		fileID := aes.DecryptID(*req.FileID)
		entity.FileID = &fileID
	}
	err = svc.repo.Insert(entity)
	err = errors.Wrap(err, "db: insert")
	return
}

func (svc *Service) GetPost(req httpPost.GetPostRequest) (resp []httpPost.GetPostResponse, err error) {
	resp = []httpPost.GetPostResponse{}
	entities, err := svc.repo.TakePost(req.Pagination)
	if err != nil {
		err = errors.Wrap(err, "db: take post")
		return
	}

	userIds := []int{}

	for _, entity := range entities {
		flag := false
		for _, userId := range userIds {
			if userId == entity.UserID {
				flag = true
				break
			}
		}

		if !flag {
			userIds = append(userIds, entity.UserID)
		}
	}

	users, err := svc.userSvc.GetUserByIDs(userIds)
	if err != nil {
		err = errors.Wrap(err, "user: get user by ids")
		return
	}

	copier.Copy(&resp, &entities)
	for i, entity := range resp {
		encID := aes.EncryptID(int(entity.ID))
		resp[i].EncID = encID
		resp[i].Username = users[int(entity.UserID)].Username
		resp[i].Email = users[int(entity.UserID)].Email
		resp[i].CreatedAt = entity.CreatedAt
	}
	return
}
