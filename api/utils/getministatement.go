package utils

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"tospay.com/WalletRestAPI/api/models"
	"tospay.com/WalletRestAPI/api/responses"
)

//GetMinistatement - performs GetMinistatement operations - fetches last 5
func GetMinistatement(tx *gorm.DB, requestCreated *models.TransactionRequests) (responses.TransactionResponse, error) {
	//first get Balance of the DebitAccount
	response := responses.TransactionResponse{}
	transaction := models.Transactions{}
	//DATETIME|TRANTYPE|TRANAMNT|DRCR"
	rows, err := tx.Debug().Model(&models.Transactions{}).Limit(5).Select("created_at,narration,amount,dr_cr").Order("id desc").Where("account_no = ?", requestCreated.DebitAccount).Rows()
	if err != nil {
		return responses.TransactionResponse{}, err
	}
	defer rows.Close()
	var minis []responses.Ministatement
	for rows.Next() {
		var mini responses.Ministatement
		tx.ScanRows(rows, &mini)
		//rows.Scan(&mini)
		minis = append(minis, mini)
	}

	response.Procode = requestCreated.Procode
	response.ResponseCode = Successful
	response.Remarks = "Ministatement Enquiry Successful"
	response.Reference = requestCreated.TxnRef
	amt, _ := strconv.ParseFloat("0.00", 64)
	response.Amount = amt
	response.Account = transaction.AccountNo
	bal, _ := strconv.ParseFloat(transaction.AvailableBal, 64)
	response.AvailableBalance = bal
	response.Ministatement = minis

	return response, nil
}
