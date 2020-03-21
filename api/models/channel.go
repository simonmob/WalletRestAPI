package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//Channel struct to hold Channel Data. The API authorizes requests using channel
type Channel struct {
	ID          uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Channel     string    `gorm:"size:255;not null;unique" json:"channel"`
	Description string    `gorm:"size:255;not null" json:"description"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//Prepare prepares channel datails
func (c *Channel) Prepare() {
	c.ID = 0
	c.Channel = html.EscapeString(strings.TrimSpace(c.Channel))
	c.Description = html.EscapeString(strings.TrimSpace(c.Description))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

//Validate validates channel data
func (c *Channel) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if c.Channel == "" {
			return errors.New("Channel Required")
		}
		if c.Description == "" {
			return errors.New("Description Required")
		}

		return nil
	case "auth":
		if c.Channel == "" {
			return errors.New("Channel Required")
		}
		return nil

	default:
		if c.Channel == "" {
			return errors.New("Channel Required")
		}
		if c.Description == "" {
			return errors.New("Description Required")
		}
		return nil
	}
}

//SaveChannel inserts in Channel Table
func (c *Channel) SaveChannel(db *gorm.DB) (*Channel, error) {

	var err error
	err = db.Debug().Create(&c).Error
	if err != nil {
		return &Channel{}, err
	}
	return c, nil
}

//FindAllChannels gets a list of all Channels
func (c *Channel) FindAllChannels(db *gorm.DB) (*[]Channel, error) {
	var err error
	channels := []Channel{}
	err = db.Debug().Model(&Channel{}).Limit(100).Find(&channels).Error
	if err != nil {
		return &[]Channel{}, err
	}
	return &channels, err
}
