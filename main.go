/**
 * Created by VoidArtanis on 10/22/2017
 */

package main

import (
	"os"

	"github.com/Artexus/api-matthew-backend/routes"
	"github.com/Artexus/api-matthew-backend/shared"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var DB = make(map[string]string)

func main() {
	godotenv.Load()

	//Db Connect and Close
	shared.Init()
	defer shared.CloseDb()

	//Init Gin
	r := gin.Default()
	routes.InitRouter(r)

	//Run Server
	r.Run(":" + os.Getenv("SERVER_PORT"))
}
