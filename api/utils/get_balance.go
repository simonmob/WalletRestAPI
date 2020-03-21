package utils

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"tospay.com/WalletRestAPI/api/models"
	"tospay.com/WalletRestAPI/api/responses"
)

//GetBalance - performs GETBALANCE operations
func GetBalance(tx *gorm.DB, requestCreated *models.TransactionRequests) (responses.TransactionResponse, error) {
	//first get Balance of the DebitAccount
	response := responses.TransactionResponse{}
	cbalance := models.Accounts{}

	err := tx.Debug().Model(&models.Accounts{}).Where("account_no = ?", requestCreated.DebitAccount).Take(&cbalance).Error
	if err != nil {
		return responses.TransactionResponse{}, err
	}
	response.Procode = requestCreated.Procode
	response.ResponseCode = Successful
	response.Remarks = "Balance Enquiry Successful"
	response.Reference = requestCreated.TxnRef
	amt, _ := strconv.ParseFloat("0.00", 64)
	response.Amount = amt
	response.Account = cbalance.AccountNo
	bal, _ := strconv.ParseFloat(cbalance.AvailableBal, 64)
	response.AvailableBalance = bal

	return response, nil
}
