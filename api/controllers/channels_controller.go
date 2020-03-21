package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tospay.com/WalletRestAPI/api/models"
	"tospay.com/WalletRestAPI/api/responses"
)

//CreateChannel receives create Channel requests
func (server *Server) CreateChannel(c *gin.Context) {

	channel := models.Channel{}
	err := c.BindJSON(&channel)
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
		return
	}
	channel.Prepare()
	err = channel.Validate("")
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
		return
	}
	channelCreated, err := channel.SaveChannel(server.DB)

	if err != nil {
		responses.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c, http.StatusCreated, channelCreated)
}

//GetChannels gets a list of all registered Channels
func (server *Server) GetChannels(c *gin.Context) {

	channel := models.Channel{}

	channels, err := channel.FindAllChannels(server.DB)
	if err != nil {
		responses.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(c, http.StatusOK, channels)
}
