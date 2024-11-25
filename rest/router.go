package rest

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
    r := gin.New()

    r.Use(cors.New(cors.Config{
        AllowOrigins:  []string{"*"},
        AllowMethods:  []string{"GET", "POST", "DELETE", "OPTIONS"},
        AllowHeaders:  []string{"Origin", "Authorization", "Content-Type"},
        ExposeHeaders: []string{"Content-Type", "Authorization"},
        MaxAge:        50 * time.Second,
    }))

    r.POST("/mailinggo/offers", SendEmailsHandler)
    r.POST("/mailinggo/openCart", PostCartNotification)
    r.POST("/mailinggo/userSubscription", HandleSubscribedUser)
    r.POST("/mailinggo/articleOffer", HandleArticleOffer)


    
    return r
}
