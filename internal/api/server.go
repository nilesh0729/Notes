package api

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/internal/db/Result"
	"github.com/nilesh0729/Notes/internal/tokens"
	"github.com/nilesh0729/Notes/internal/util"
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

	// Configure CORS to allow requests from frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins, can be restricted to specific domain
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/user", server.CreateUser)
	router.POST("/login", server.LoginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/notes", server.CreateNote)
	authRoutes.GET("/notes/:id", server.GetNoteById)
	authRoutes.GET("/notes", server.ListNotes)
	authRoutes.PUT("/notes/:id", server.UpdateNote)
	authRoutes.DELETE("/notes/:id", server.DeleteNote)

	authRoutes.POST("/tags", server.CreateTags)
	authRoutes.GET("/tags/:id", server.GetTag)
	authRoutes.GET("/tags", server.ListTags)
	authRoutes.DELETE("/tags/:id", server.DeleteTag)
	
	authRoutes.GET("/tags/:id/notes", server.ListNotesForTag)

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