package auth

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/Artexus/api-matthew-backend/constant"
	httpToken "github.com/Artexus/api-matthew-backend/http/token"
	httpUser "github.com/Artexus/api-matthew-backend/http/user"
	model "github.com/Artexus/api-matthew-backend/models/user"
	"github.com/Artexus/api-matthew-backend/repositories/user"
	"github.com/Artexus/api-matthew-backend/services/token"
	"github.com/Artexus/api-matthew-backend/utils/hash"
	"github.com/Artexus/api-matthew-backend/utils/jwt"
)

type Service struct {
	svcToken token.Service
	repo     user.RepositoryInterface
}

var instance *Service

func NewService(svcToken token.Service, repo user.RepositoryInterface) *Service {
	if instance == nil {
		instance = &Service{
			svcToken: svcToken,
			repo:     repo,
		}
	}
	return instance
}

func (svc *Service) Login(identifier, password string) (resp httpUser.LoginResponse, err error) {
	userData, err := svc.repo.TakeUserByIdentifier(strings.ToLower(identifier))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = constant.ErrUserNotExist
		}
		return
	}

	if userData.Password != hash.GenerateHashToken(password) {
		err = constant.ErrWrongPassword
		return
	}
	requests := make([]httpToken.InsertTokenRequest, 0)
	accessTokenExpired := time.Now().Add(constant.AccessTokenInterval).Unix()
	accessToken, err := jwt.GenerateToken(userData.EncID(), userData.Email, userData.Username, accessTokenExpired, constant.AccessTokenSignedKey)
	if err != nil {
		err = constant.ErrGenerateTokenFailed
		return
	}
	requests = append(requests, httpToken.InsertTokenRequest{
		Type:      constant.AccessTokenType,
		UserID:    int(userData.ID),
		Token:     accessToken,
		ExpiredAt: accessTokenExpired,
	})

	refreshTokenExpired := time.Now().Add(constant.RefreshTokenInterval).Unix()
	refreshToken, err := jwt.GenerateToken(userData.EncID(), userData.Email, userData.Username, refreshTokenExpired, constant.RefreshTokenSignedKey)
	if err != nil {
		err = constant.ErrGenerateTokenFailed
		return
	}
	requests = append(requests, httpToken.InsertTokenRequest{
		Type:      constant.RefreshTokenType,
		UserID:    int(userData.ID),
		Token:     refreshToken,
		ExpiredAt: refreshTokenExpired,
	})

	err = svc.svcToken.UpsertAccessAndRefreshToken(requests)
	if err != nil {
		err = errors.Wrap(err, "upsert access and refresh token")
		return
	}
	resp.ID = userData.EncID()
	resp.Email = userData.Email
	resp.Username = userData.Username
	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken
	return
}

func (svc *Service) Register(req httpUser.RegisterRequest) (err error) {
	exist, err := svc.repo.ExistUsername(req.Username)
	if err != nil {
		err = errors.Wrap(err, "db: exist username")
		return
	} else if exist {
		err = constant.ErrExistUsername
		return
	}

	exist, err = svc.repo.ExistEmail(req.Email)
	if err != nil {
		err = errors.Wrap(err, "db: exist email")
		return
	} else if exist {
		err = constant.ErrExistEmail
		return
	}

	entity := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hash.GenerateHashToken(req.Password),
	}

	return errors.Wrap(svc.repo.Create(entity), "[ERROR] user insert")
}
