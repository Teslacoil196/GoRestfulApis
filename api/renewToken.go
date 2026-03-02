package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewUserTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewUserTokenRespons struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpireTime time.Time `json:"access_token_expire_time"`
}

func (server *Server) renewUserToken(ctx *gin.Context) {
	var request renewUserTokenRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorHandler(err))
		return
	}

	session, err := server.db.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorHandler(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		ctx.JSON(http.StatusNotFound, errorHandler(err))
		return
	}
	if session.Username != payload.Username {
		err := fmt.Errorf("Incorrect session")
		ctx.JSON(http.StatusNotFound, errorHandler(err))
		return
	}
	if session.RefreshToken != request.RefreshToken {
		err := fmt.Errorf("Incorrect refresh token")
		ctx.JSON(http.StatusNotFound, errorHandler(err))
		return
	}
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("Expried session")
		ctx.JSON(http.StatusNotFound, errorHandler(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(payload.Username, server.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	reply := renewUserTokenRespons{
		AccessToken:           accessToken,
		AccessTokenExpireTime: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, reply)

}
