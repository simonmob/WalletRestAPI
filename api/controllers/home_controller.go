package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tospay.com/WalletRestAPI/api/responses"
)

//Home return welcome message to the API consumers
func (server *Server) Home(c *gin.Context) {
	responses.JSON(c, http.StatusOK, "Welcome To Tospay Wallet API")
}
