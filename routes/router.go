/**
 * Created by VoidArtanis on 10/22/2017
 */

package routes

import (
	authcontroller "github.com/Artexus/api-matthew-backend/controllers/auth"
	postcontroller "github.com/Artexus/api-matthew-backend/controllers/post"
	"github.com/Artexus/api-matthew-backend/middlewares"
	postrepo "github.com/Artexus/api-matthew-backend/repositories/post"
	tokenrepo "github.com/Artexus/api-matthew-backend/repositories/token"
	userrepo "github.com/Artexus/api-matthew-backend/repositories/user"
	authservices "github.com/Artexus/api-matthew-backend/services/auth"
	postservices "github.com/Artexus/api-matthew-backend/services/post"
	tokenservices "github.com/Artexus/api-matthew-backend/services/token"
	userservices "github.com/Artexus/api-matthew-backend/services/user"
	"github.com/Artexus/api-matthew-backend/shared"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {
	db := shared.GetDb()
	tokenRepo := tokenrepo.NewRepository(db)
	authRepo := userrepo.NewRepository(db)
	postRepo := postrepo.NewRepository(db)

	tokenSvc := tokenservices.NewService(tokenRepo)
	authService := authservices.NewService(*tokenSvc, authRepo)
	userSvc := userservices.NewService(authRepo)
	postSvc := postservices.NewService(postRepo, *userSvc)

	authController := authcontroller.NewController(*authService, *tokenSvc)
	postController := postcontroller.NewController(*postSvc)
	engine.Use(middlewares.CORSMiddleware())

	auth := engine.Group("auth")
	registerAuthRoute(auth, authController)

	post := engine.Group("posts", middlewares.AuthHandler())
	registerPostRoute(post, postController)
}

func registerAuthRoute(router *gin.RouterGroup, handler *authcontroller.Controller) {
	router.POST("login", handler.Login)
	router.POST("register", handler.Register)
	router.POST("token", handler.GenerateToken)
}

func registerPostRoute(router *gin.RouterGroup, handler *postcontroller.Controller) {
	router.GET("", handler.GetPost)
	router.POST("", handler.InsertPost)
}
