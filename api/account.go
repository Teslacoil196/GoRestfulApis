package api

import (
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/token"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountServerParams struct {
	//Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var request CreateAccountServerParams

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  int64(0),
		Currency: request.Currency,
	}

	account, err := server.db.CreateAccount(ctx, arg)
	if err != nil {

		if pgerr, ok := err.(*pq.Error); ok {
			switch pgerr.Code.Name() {
			case "foregin_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorHandler(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type GetAccountParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var request GetAccountParam

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	account, err := server.db.GetAccount(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorHandler(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("Account doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorHandler(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountsParam struct {
	PageID   int32 `form:"page_id" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var request ListAccountsParam

	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}

	accounts, err := server.db.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

type DeleteAccountParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var request DeleteAccountParam

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	err = server.db.DeleteAccount(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorHandler(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Deleted account with id %d", request.ID),
		"id":      request.ID,
	})
}

type UpdateAccountServerParams struct {
	ID      int64 `json:"id" binding:"required"`
	Balance int64 `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var request UpdateAccountServerParams

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	account, err := server.db.GetAccount(ctx, request.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      request.ID,
		Balance: account.Balance + request.Balance,
	}

	accountUpdated, err := server.db.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}
	ctx.JSON(http.StatusOK, accountUpdated)
}
