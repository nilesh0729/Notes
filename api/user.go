package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	Database "github.com/nilesh0729/Notes/db/Result"
	"github.com/nilesh0729/Notes/util"
)

type UserResponseFormat struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func UserResponse(user Database.User) UserResponseFormat {
	return UserResponseFormat{
		Username: user.Username,
		Password: "********",
		Email:    "******@gmail.com",
	}
}

type CreateUserRequest struct {
	Username string `json:"Username" binding:"required,alphanum"`
	Password string `json:"Password" binding:"required,min=8"`
	Email    string `json:"Email" binding:"required,email"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err!= nil{
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	arg := Database.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, UserResponse(user))
}

type GetUserRequest struct{
	Username string `uri:"username" binding:"required,alphanum,min=6"`
}

func (server *Server) Getuser(ctx *gin.Context){
	var req GetUserRequest

	err := ctx.BindUri(&req)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, UserResponse(user))
}
