package handler

import (
	"github.com/PSNAppz/Fold-ELK/db"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	DB       db.Database
	Logger   zerolog.Logger
	ESClient *elasticsearch.Client
}

func New(database db.Database, esClient *elasticsearch.Client, logger zerolog.Logger) *Handler {
	return &Handler{
		DB:       database,
		ESClient: esClient,
		Logger:   logger,
	}
}

// All the routes are defined here
func (h *Handler) Register(group *gin.RouterGroup) {

	//Define routes for users
	group.POST("/users", h.CreateUser)
	group.PATCH("/users/:id", h.UpdateUser)
	group.GET("/users/:id", h.GetUser)
	group.DELETE("/users/:id", h.DeleteUser)
	group.GET("/users", h.GetUsers)

	// Define routes for hashtags
	group.POST("/hashtags", h.CreateHashtag)
	group.PATCH("/hashtags/:id", h.UpdateHashtag)
	group.GET("/hashtags/:id", h.GetHashtag)
	group.DELETE("/hashtags/:id", h.DeleteHashtag)
	group.GET("/hashtags", h.GetHashtags)

	// Define routes for posts
	group.POST("/projects", h.CreateProject)
	group.PATCH("/projects/:id", h.UpdateProject)
	group.GET("/projects/:id", h.GetProject)
	group.DELETE("/projects/:id", h.DeleteProject)
	group.GET("/projects", h.GetProjects)

	// Define routes for search
	group.GET("/search", h.SearchProjects)
	group.GET("/fuzzy_search/", h.FuzzySearchProjects)
}
