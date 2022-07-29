package post

import (
	"log"
	"net/http"

	"github.com/Artexus/api-matthew-backend/constant"
	httpPost "github.com/Artexus/api-matthew-backend/http/post"
	"github.com/Artexus/api-matthew-backend/services/post"
	"github.com/Artexus/api-matthew-backend/utils/jwt"
	"github.com/Artexus/api-matthew-backend/utils/rest"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc post.Service
}

var instance *Controller

func NewController(svc post.Service) *Controller {
	if instance == nil {
		instance = &Controller{
			svc: svc,
		}
	}
	return instance
}

func (ctrl *Controller) InsertPost(ctx *gin.Context) {
	req := httpPost.InsertPostRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, nil, &map[string]string{
			"body": constant.ErrInvalid.Error(),
		})
		return
	}

	req.UserID, _ = jwt.ExtractIDToken(ctx.GetHeader("Authorization"), constant.AccessTokenSignedKey)
	err = ctrl.svc.InsertPost(req)
	if err != nil {
		log.Println("[ERROR] insert post: ", err.Error())
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil, nil)
		return
	}
	rest.ResponseOutput(ctx, http.StatusOK, nil, nil)
}

func (ctrl *Controller) GetPost(ctx *gin.Context) {
	req := httpPost.GetPostRequest{}
	err := ctx.BindQuery(&req)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, nil, &map[string]string{
			"body": constant.ErrInvalid.Error(),
		})
		return
	}

	req.Pagination.Paginate()
	resp, err := ctrl.svc.GetPost(req)
	if err != nil {
		log.Println("[ERROR] get post: ", err.Error())
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil, nil)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, resp)
}
