package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

//Accounts hold Customer account details
type Accounts struct {
	gorm.Model
	AccountNo    string `gorm:"size:255;not null;unique" json:"account_no"`
	Dormant      string `gorm:"size:255" json:"dormant"`
	Approved     bool
	ActualBal    string `gorm:"type:numeric(19,2)" json:"actual_bal"`
	AvailableBal string `gorm:"type:numeric(19,2)" json:"available_bal"`
}

//Prepare prepares Account Data.
func (a *Accounts) Prepare() {
	a.AccountNo = html.EscapeString(strings.TrimSpace(a.AccountNo))
	a.Dormant = "N"
	a.Approved = true
	a.ActualBal = "0.00"
	a.AvailableBal = "0.00"
}

//Validate validates account data
func (a *Accounts) Validate() error {

	if a.AccountNo == "" {
		return errors.New("AccountNo Required")
	}

	if a.AccountNo == "" {
		return errors.New("AccountNo Required")
	}
	if len(a.AccountNo) != 12 {
		return errors.New("PhoneNumber length should be 12")
	}

	return nil
}

//SaveAccount inserts account data in Accounts table
func (a *Accounts) SaveAccount(db *gorm.DB) (*Accounts, error) {
	var err error
	err = db.Debug().Model(&Accounts{}).Create(&a).Error
	if err != nil {
		return &Accounts{}, err
	}
	return a, nil
}

//UpdateAccount updates AccountsTable with the new balances after processing
func (a *Accounts) UpdateAccount(db *gorm.DB, account string) (*Accounts, error) {

	err := db.Debug().Model(&Accounts{}).Where("account_no = ?", account).Updates(Accounts{AvailableBal: a.AvailableBal, ActualBal: a.ActualBal}).Error
	if err != nil {
		return &Accounts{}, err
	}

	// This is the display the updated record
	err = db.Debug().Model(&Accounts{}).Where("account_no = ?", account).Take(&a).Error
	if err != nil {
		return &Accounts{}, err
	}
	return a, nil
}
