package token

import (
	"time"

	"github.com/Artexus/api-matthew-backend/constant"
	httpToken "github.com/Artexus/api-matthew-backend/http/token"
	dbToken "github.com/Artexus/api-matthew-backend/models/token"
	"github.com/Artexus/api-matthew-backend/repositories/token"
	"github.com/Artexus/api-matthew-backend/utils/jwt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Service struct {
	repo token.RepositoryInterface
}

var instance *Service

func NewService(repo token.RepositoryInterface) *Service {
	if instance == nil {
		instance = &Service{
			repo: repo,
		}
	}
	return instance
}

func (svc *Service) UpsertAccessAndRefreshToken(requests []httpToken.InsertTokenRequest) (err error) {
	var exist bool
	for _, req := range requests {
		entity := dbToken.Token{
			UserID:    req.UserID,
			Type:      req.Type,
			ExpiredAt: req.ExpiredAt,
			Token:     req.Token,
		}
		exist, err = svc.repo.ExistByUserID(entity.UserID, entity.Type)
		if err != nil {
			err = errors.Wrap(err, "db: exist by user id")
			return
		}

		if !exist {
			err = svc.repo.Create(entity)
			if err != nil {
				err = errors.Wrap(err, "db: create")
				return
			}
		} else {
			err = svc.repo.Update(entity)
			if err != nil {
				err = errors.Wrap(err, "db: update")
				return
			}
		}
	}
	return
}

func (svc *Service) GenerateToken(req httpToken.GenerateTokenRequest) (resp httpToken.GenerateTokenResponse, err error) {
	signedKey := constant.AccessTokenSignedKey
	if req.TokenType != constant.AccessTokenType {
		signedKey = constant.RefreshTokenSignedKey
	}

	userData, err := jwt.ExtractToken(req.Token, signedKey)
	if err != nil {
		err = errors.Wrap(err, "extract id token")
		return
	}

	tokenData, err := svc.repo.TakeByToken(req.Token, req.TokenType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = constant.ErrTokenNotFound
			return
		}
		err = errors.Wrap(err, "db: take by token")
		return
	}

	if tokenData.UserID != userData.UserID {
		err = constant.ErrInvalidID
		return
	}

	if time.Now().After(time.Unix(tokenData.ExpiredAt, 0)) {
		err = constant.ErrTokenExpired
		return
	}

	reqTokenData, err := svc.repo.TakeByUser(tokenData.UserID, req.RequestType)
	if err != nil {
		err = errors.Wrap(err, "db: take by token")
		return
	}

	if time.Now().After(time.Unix(reqTokenData.ExpiredAt, 0)) {
		resp.Token = reqTokenData.Token
		resp.Type = req.RequestType
		return
	}

	var request httpToken.InsertTokenRequest
	var token string
	if req.RequestType == constant.AccessTokenType {
		accessTokenExpired := time.Now().Add(constant.AccessTokenInterval).Unix()
		token, err = jwt.GenerateToken(userData.EncID, userData.Email, userData.Username, accessTokenExpired, constant.AccessTokenSignedKey)
		if err != nil {
			err = errors.Wrap(err, "generate access token")
			return
		}
		request.ExpiredAt = accessTokenExpired
		request.Type = req.RequestType
	} else {
		refreshTokenExpired := time.Now().Add(constant.RefreshTokenInterval).Unix()
		token, err = jwt.GenerateToken(userData.EncID, userData.Email, userData.Username, refreshTokenExpired, constant.RefreshTokenSignedKey)
		if err != nil {
			err = errors.Wrap(err, "generate refresh token")
			return
		}
		request.ExpiredAt = refreshTokenExpired
		request.Type = req.RequestType
	}

	request.UserID = userData.UserID
	request.Token = token

	err = svc.UpsertAccessAndRefreshToken([]httpToken.InsertTokenRequest{request})
	if err != nil {
		err = errors.Wrap(err, "upsert access and refresh token")
		return
	}

	resp.Token = request.Token
	resp.Type = request.Type
	return
}
