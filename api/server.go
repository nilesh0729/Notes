package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/db/Result"
	"github.com/nilesh0729/Notes/tokens"
	"github.com/nilesh0729/Notes/util"
)

type Server struct {
	config     util.Config
	store      Database.Store
	tokenMaker tokens.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store Database.Store) (*Server, error) {
	tokenMaker, err := tokens.NewPasetoMaker(config.Secret)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	router.POST("/user", server.CreateUser)
	router.POST("/login", server.LoginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/notes", server.CreateNote)
	authRoutes.GET("/notes/:id", server.GetNoteById)
	authRoutes.GET("/notes", server.ListNotes)

	authRoutes.POST("/tags", server.CreateTags)
	authRoutes.GET("/tags/:id", server.GetTag)
	authRoutes.GET("/tags", server.ListTags)

	authRoutes.POST("/note_tags", server.AddTagToNote)

	server.router = router

	return server, nil
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func errResponse (err error)gin.H{
	return gin.H{"error" : err.Error() }
}