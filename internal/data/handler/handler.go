package handler

import (
	"pet-app/internal/data"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	searchURLs map[int]data.SearchURL
	count      *int
	tokenAdmin string
}

func New(inputSearchMap map[int]data.SearchURL, inputCount *int, inputToken string) Handler {
	return Handler{
		searchURLs: inputSearchMap,
		count:      inputCount,
		tokenAdmin: inputToken,
	}
}

func (h Handler) Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", h.showData)
	admin := router.Group("/admin", h.authorization(h.tokenAdmin))
	{
		admin.POST("/", h.addSearchURL)
		admin.DELETE("/", h.deleteSearchURL)
	}
	return router
}
