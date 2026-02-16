package api

import (
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/token"
	"TeslaCoil196/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	db         db.Store
	config     util.Config
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(db db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPastoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Could not create tokenMaker %s", errorHandler(err))
	}
	server := &Server{
		config:     config,
		db:         db,
		tokenMaker: tokenMaker,
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validateCurrency)
	}
	server.setUpRouter()

	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.GET("/", HelloTheir)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoot := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoot.POST("/account", server.createAccount)
	authRoot.GET("/account/:id", server.getAccount)
	authRoot.GET("/account", server.listAccounts)
	authRoot.DELETE("/account/delete/:id", server.deleteAccount)
	authRoot.POST("/account/update", server.updateAccount)

	authRoot.POST("/transfer", server.createTransfer)

	server.router = router
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
