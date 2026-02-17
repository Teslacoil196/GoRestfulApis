package api

import (
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/token"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateTransferServerParams struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var request CreateTransferServerParams

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	fromAccount, valid := server.validateAccounts(ctx, request.FromAccountID, request.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != fromAccount.Owner {
		err := errors.New("Account doesn't belong to authencated user")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorHandler(err))
	}

	_, valid = server.validateAccounts(ctx, request.ToAccountID, request.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		Amount:        request.Amount,
	}

	result, err := server.db.TranferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateAccounts(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {

	account, err := server.db.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorHandler(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("Currency mismatch for account %d; %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return account, false
	}

	return account, true
}
