package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

//Transactions struct to hold completed Transactions details
type Transactions struct {
	gorm.Model
	MsgType      string `gorm:"size:6;not null"`
	Procode      string `gorm:"size:6;not null"`
	Channel      string `gorm:"size:20;not null"`
	TxnRef       string `gorm:"size:12;not null"`
	AccountNo    string `gorm:"size:12;not null"`
	Amount       string `gorm:"type:numeric(19,2);not null"`
	DrCr         string `gorm:"size:2;not null"`
	Narration    string `gorm:"size:200;not null"`
	AvailableBal string `gorm:"type:numeric(19,2);not null"`
}

//Prepare prepares Transactions data before insert
func (t *Transactions) Prepare() {
	t.MsgType = html.EscapeString(strings.TrimSpace(t.MsgType))
	t.Procode = html.EscapeString(strings.TrimSpace(t.Procode))
	t.Channel = html.EscapeString(strings.TrimSpace(t.Channel))
	t.TxnRef = html.EscapeString(strings.TrimSpace(t.TxnRef))
	t.Amount = html.EscapeString(strings.TrimSpace(t.Amount))
	t.AccountNo = html.EscapeString(strings.TrimSpace(t.AccountNo))
	t.Narration = html.EscapeString(strings.TrimSpace(t.Narration))
	t.DrCr = html.EscapeString(strings.TrimSpace(t.DrCr))
	t.AvailableBal = html.EscapeString(strings.TrimSpace(t.AvailableBal))
}

//Validate validtaes Transactions data before inserts
func (t *Transactions) Validate() error {

	if t.MsgType == "" {
		return errors.New("MsgType Required")
	}

	if t.Procode == "" {
		return errors.New("Procode Required")
	}
	if t.Channel == "" {
		return errors.New("Channel Required")
	}
	if t.TxnRef == "" {
		return errors.New("TxnRef Required")
	}
	if t.Amount == "" {
		return errors.New("Amount Required")
	}
	if t.Narration == "" {
		return errors.New("Narration Required")
	}
	if len(t.AccountNo) != 12 {
		return errors.New("AccountNo length should be 12")
	}
	if len(t.DrCr) != 12 {
		return errors.New("DRCR Required")
	}

	return nil
}

//SaveTransactions inserts transaction details in Transactions table
func (t *Transactions) SaveTransactions(db *gorm.DB) (*Transactions, error) {
	var err error
	err = db.Debug().Model(&TransactionRequests{}).Create(&t).Error
	if err != nil {
		return &Transactions{}, err
	}
	return t, nil
}
