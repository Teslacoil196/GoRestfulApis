package api

import (
	db "TeslaCoil196/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	db     db.Store
	router *gin.Engine
}

func NewServer(db db.Store) *Server {
	server := &Server{db: db}
	router := gin.Default()

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validateCurrency)
	}
	router.POST("/users", server.createUser)

	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/account", server.listAccounts)
	router.DELETE("/account/delete/:id", server.deleteAccount)
	router.POST("/account/update", server.updateAccount)

	router.POST("/transfer", server.createTransfer)

	router.GET("/", HelloTheir)

	server.router = router
	return server
}

func HelloTheir(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hellow their !")
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

func errorHandler(err error) gin.H {
	return gin.H{"error": err.Error()}
}
