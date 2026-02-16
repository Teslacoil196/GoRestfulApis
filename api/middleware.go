package api

import (
	"TeslaCoil196/token"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKet  = "authorization"
	authorizationType       = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKet)
		if len(authHeader) == 0 {
			err := errors.New("Authorization header not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorHandler(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("Invalid authorization header provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorHandler(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authorizationType != authType {
			err := errors.New("Unsupported AuthType ")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorHandler(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()

	}
}
