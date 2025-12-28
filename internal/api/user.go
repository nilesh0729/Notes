package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	Database "github.com/nilesh0729/Notes/internal/db/Result"
	"github.com/nilesh0729/Notes/internal/util"
)

type UserResponseFormat struct {
	Username string `json:"username" binding:"required,alphanum,min=6"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
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

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=6"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserResponse struct{
	AccessToken string `json:"access_token"`
	User UserResponseFormat `json:"user"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	res := LoginUserResponse{
		AccessToken: accessToken,
		User:        UserResponse(user),
	}

	ctx.JSON(http.StatusOK, res)
}