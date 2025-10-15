package api

import (
	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/db/Result"
)

type Server struct {
	store  Database.Store
	router *gin.Engine
}

func NewServer(store Database.Store) *Server{
	server := &Server{
		store : store,
	}
	router := gin.Default()

	router.POST("/notes", server.CreateNote)
	
	router.GET("/notes/:id", server.GetNoteById)

	router.GET("/notes", server.ListNotes)

	router.POST("/tags", server.CreateTags)

	router.GET("/tags/:id", server.GetTag)

	router.GET("/tags", server.ListTags)

	router.POST("/note_tags", server.AddTagToNote)

	
	server.router = router

	return server
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func errResponse (err error)gin.H{
	return gin.H{"error" : err.Error() }
}