package api

import (
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/util"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CreateUserServerParams struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username      string    `json:"username"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	CreatedAt     time.Time `json:"created_at"`
	LastPassReset time.Time `json:"last_pass_reset"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var request CreateUserServerParams

	//fmt.Println("CreateUser called")

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	hashedPassword, err := util.HashedPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       request.Username,
		HashedPassword: hashedPassword,
		FullName:       request.FullName,
		Email:          request.Email,
	}
	//fmt.Println("about to make the call")
	user, err := server.db.CreateUser(ctx, arg)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			switch pgerr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorHandler(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	response := newUserResponse(user)
	//fmt.Println(" status ok ")
	ctx.JSON(http.StatusOK, response)
}

func newUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:      user.Username,
		FullName:      user.FullName,
		Email:         user.Email,
		CreatedAt:     user.CreatedAt,
		LastPassReset: user.LastPassReset,
	}
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID              uuid.UUID    `json:"session_id"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpireTime  time.Time    `json:"access_token_expire_time"`
	User                   UserResponse `json:"user"`
	RefreshToken           string       `json:"refres_token"`
	RefreshTokenExpireTime time.Time    `json:"refresh_token_expire_time"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var request loginUserRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	user, err := server.db.GetUser(ctx, request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorHandler(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}
	err = util.CheckPasswords(request.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorHandler(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	session, err := server.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    " ", // TODO add dis somewhere down the line
		ClientIp:     " ",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	reply := loginUserResponse{
		SessionID:              session.ID,
		AccessToken:            accessToken,
		AccessTokenExpireTime:  accessPayload.ExpiredAt,
		User:                   newUserResponse(user),
		RefreshToken:           refreshToken,
		RefreshTokenExpireTime: refreshPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, reply)

}
