package app

/*
	application will interact with our http server
*/
import (
	"github.com/gin-gonic/gin"
	"github.com/tamihyo/bookstore_users-api/logger"
)

/*
every request in application will be handled by this router
*/
var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	logger.Info("about to start the application")
	router.Run(":8085")
}
