package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

//TransactionRequests struct to hold all transaction requests incoming from different channels
type TransactionRequests struct {
	gorm.Model
	MsgType       string `gorm:"size:4;not null" json:"msg_type"`
	Procode       string `gorm:"size:6;not null" json:"pro_code"`
	Channel       string `gorm:"size:20;not null" json:"channel"`
	TxnRef        string `gorm:"size:12;not null;unique" json:"txn_ref"`
	Amount        string `gorm:"type:numeric(19,2);not null" json:"amount"`
	ResponseCode  string `gorm:"size:6" json:"ResponseCode"`
	Remarks       string `gorm:"size:200" json:"remarks"`
	Narration     string `gorm:"size:200;not null" json:"narration"`
	DebitAccount  string `gorm:"size:12;not null" json:"debit_account"`
	CreditAccount string `gorm:"size:12" json:"credit_account"`
}

//Prepare prepares TransactionRequests data
func (r *TransactionRequests) Prepare() {
	r.MsgType = html.EscapeString(strings.TrimSpace(r.MsgType))
	r.Procode = html.EscapeString(strings.TrimSpace(r.Procode))
	r.Channel = html.EscapeString(strings.TrimSpace(r.Channel))
	r.TxnRef = html.EscapeString(strings.TrimSpace(r.TxnRef))
	r.Amount = html.EscapeString(strings.TrimSpace(r.Amount))
	r.Narration = html.EscapeString(strings.TrimSpace(r.Narration))
	r.DebitAccount = html.EscapeString(strings.TrimSpace(r.DebitAccount))
	r.CreditAccount = html.EscapeString(strings.TrimSpace(r.CreditAccount))
}

//Validate validates TransactionRequests data.
func (r *TransactionRequests) Validate() error {

	if r.MsgType == "" {
		return errors.New("MsgType Required")
	}

	if r.Procode == "" {
		return errors.New("Procode Required")
	}
	if r.Channel == "" {
		return errors.New("Channel Required")
	}
	if r.TxnRef == "" {
		return errors.New("TxnRef Required")
	}
	if r.Amount == "" {
		return errors.New("Amount Required")
	}
	if r.Narration == "" {
		return errors.New("Narration Required")
	}
	if len(r.DebitAccount) != 12 {
		return errors.New("AccountNo length should be 12")
	}

	return nil
}

//SaveTransactionRequest inserts incoming TransactionRequests into TransactionRequests table
func (r *TransactionRequests) SaveTransactionRequest(db *gorm.DB) (*TransactionRequests, error) {
	var err error
	err = db.Debug().Model(&TransactionRequests{}).Create(&r).Error
	if err != nil {
		return &TransactionRequests{}, err
	}
	return r, nil
}

//UpdateTransactionResponse updates status of the transaction after processing
func (r *TransactionRequests) UpdateTransactionResponse(db *gorm.DB, ref string) (*TransactionRequests, error) {
	// Read the data posted and Start processing the request data
	err := db.Debug().Model(&TransactionRequests{}).Where("txn_ref = ?", ref).Updates(TransactionRequests{ResponseCode: r.ResponseCode, Remarks: r.Remarks}).Error
	if err != nil {
		return &TransactionRequests{}, err
	}

	// This is the display the updated record
	err = db.Debug().Model(&TransactionRequests{}).Where("txn_ref = ?", ref).Take(&r).Error
	if err != nil {
		return &TransactionRequests{}, err
	}
	return r, nil

}
