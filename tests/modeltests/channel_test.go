package modeltests

import (
	"log"
	"reflect"
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"tospay.com/WalletRestAPI/api/models"
)

func TestChannel_SaveChannel(t *testing.T) {

	err := refreshChannelTable()
	if err != nil {
		log.Fatal(err)
	}
	newChannel := models.Channel{
		ID:          1,
		Channel:     "USSD",
		Description: "test Channel",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	savedChannel, err := newChannel.SaveChannel(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the channel: %v\n", err)
		return
	}

	if !reflect.DeepEqual(savedChannel.Channel, newChannel.Channel) {
		t.Errorf("Channel.SaveChannel() = %v, want %v", savedChannel, newChannel)
	}
	if !reflect.DeepEqual(savedChannel.Description, newChannel.Description) {
		t.Errorf("Channel.SaveChannel() = %v, want %v", savedChannel, newChannel)
	}

	// assert.Equal(t, newUser.ID, savedUser.ID)
	// assert.Equal(t, newUser.Email, savedUser.Email)
	// assert.Equal(t, newUser.Nickname, savedUser.Nickname)

}

func TestChannel_FindAllChannels(t *testing.T) {

	_, err := loadOneChannel()
	if err != nil {
		log.Fatal(err)
	}

	channels, err := channelInstance.FindAllChannels(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the channels: %v\n", err)
		return
	}
	//assert.Equal(t, len(*users), 2)

	for _, v := range *channels {
		if !reflect.DeepEqual("USSD", v.Channel) {
			t.Errorf("Channel.FindAllChannels() got = %v, want %v", v.Channel, "USSD")
		}
	}

}
