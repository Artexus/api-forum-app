package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Detail  interface{} `json:"detail,omitempty"`
}

func ResponseOutput(ctx *gin.Context, status int, message *string, detail *map[string]string) {
	response := Response{
		Status: http.StatusText(status),
	}
	if detail != nil {
		response.Detail = detail
	} else if message != nil {
		response.Message = *message
	}
	log.Println(response)
	ctx.JSON(status, response)

}

func ResponseData(ctx *gin.Context, status int, entity interface{}) {
	response := Response{
		Status: http.StatusText(status),
		Detail: entity,
	}
	ctx.JSON(status, response)

}
