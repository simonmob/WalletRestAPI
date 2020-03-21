package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tospay.com/WalletRestAPI/api/authorization"
	"tospay.com/WalletRestAPI/api/models"
	"tospay.com/WalletRestAPI/api/responses"
)

//ChannelAuth works as channel login to get token for use in future requests
func (server *Server) ChannelAuth(c *gin.Context) {
	channel := models.Channel{}
	err := c.BindJSON(&channel)
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
		return
	}

	channel.Prepare()
	err = channel.Validate("auth")
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.GetChannelToken(channel.Channel)
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
		return
	}
	responses.JSON(c, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

//GetChannelToken creates the channel authorization token
func (server *Server) GetChannelToken(channelstr string) (string, error) {

	var err error

	channel := models.Channel{}

	err = server.DB.Debug().Model(models.Channel{}).Where("channel = ?", channelstr).Take(&channel).Error
	if err != nil {
		return "", err
	}

	return authorization.CreateToken(channel.Channel)
}
