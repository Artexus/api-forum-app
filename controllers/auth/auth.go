package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/Artexus/api-matthew-backend/constant"
	httpToken "github.com/Artexus/api-matthew-backend/http/token"
	httpUser "github.com/Artexus/api-matthew-backend/http/user"
	"github.com/Artexus/api-matthew-backend/services/auth"
	"github.com/Artexus/api-matthew-backend/services/token"
	"github.com/Artexus/api-matthew-backend/utils/rest"
	"github.com/Artexus/api-matthew-backend/utils/validate"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc      auth.Service
	tokenSvc token.Service
}

var instance *Controller

func NewController(svc auth.Service, tokenSvc token.Service) *Controller {
	if instance == nil {
		instance = &Controller{
			svc:      svc,
			tokenSvc: tokenSvc,
		}
	}
	return instance
}

func (ctrl *Controller) Login(ctx *gin.Context) {
	req := httpUser.LoginRequest{}
	if err := ctx.Bind(&req); err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, nil, &map[string]string{
			"body": constant.ErrInvalid.Error(),
		})
		return
	}
	resp, err := ctrl.svc.Login(req.Identifier, req.Password)
	if err != nil {
		msgErr := err.Error()
		if errors.Is(err, constant.ErrUserNotExist) {
			rest.ResponseOutput(ctx, http.StatusNotFound, &msgErr, nil)
			return
		} else if errors.Is(err, constant.ErrWrongPassword) {
			rest.ResponseOutput(ctx, http.StatusNotFound, &msgErr, nil)
			return
		}

		log.Println("[ERROR] login: ", msgErr)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil, nil)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, resp)
}

func (ctrl *Controller) Register(ctx *gin.Context) {
	req := httpUser.RegisterRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, nil, &map[string]string{
			"body": constant.ErrInvalid.Error(),
		})
		return
	}

	if err = validate.ValidateRegister(req.Username, req.Email, req.Password); err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, nil, &map[string]string{
			"body": err.Error(),
		})
		return
	}

	err = ctrl.svc.Register(req)
	if errors.Is(err, constant.ErrExistEmail) || errors.Is(err, constant.ErrExistUsername) {
		msgErr := err.Error()
		rest.ResponseOutput(ctx, http.StatusConflict, &msgErr, nil)
		return
	} else if err != nil {
		log.Println("[ERROR] register: ", err.Error())
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil, nil)
		return
	}

	rest.ResponseOutput(ctx, http.StatusCreated, nil, nil)
}

func (ctrl *Controller) GenerateToken(ctx *gin.Context) {
	req := httpToken.GenerateTokenRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, nil, &map[string]string{
			"body": constant.ErrInvalid.Error(),
		})
		return
	}
	log.Println(req)
	if req.RequestType != constant.AccessTokenType && req.RequestType != constant.RefreshTokenType {
		rest.ResponseOutput(ctx, http.StatusBadRequest, nil, &map[string]string{
			"type": constant.ErrInvalid.Error(),
		})
		return
	}

	req.TokenType = constant.AccessTokenType
	if req.RequestType == constant.AccessTokenType {
		req.TokenType = constant.RefreshTokenType
	}

	resp, err := ctrl.tokenSvc.GenerateToken(req)
	if errors.Is(err, constant.ErrInvalidID) {
		rest.ResponseOutput(ctx, http.StatusBadRequest, nil, &map[string]string{
			"id": constant.ErrInvalidID.Error(),
		})
		return
	} else if errors.Is(err, constant.ErrTokenNotFound) ||
		errors.Is(err, constant.ErrTokenInvalid) ||
		errors.Is(err, constant.ErrTokenExpired) {
		rest.ResponseOutput(ctx, http.StatusNotFound, nil, &map[string]string{
			"token": err.Error(),
		})
		return
	} else if err != nil {
		log.Println("[ERROR] generate token: ", err.Error())
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil, nil)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, resp)
}
