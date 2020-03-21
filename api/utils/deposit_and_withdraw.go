package utils

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	"tospay.com/WalletRestAPI/api/models"
	"tospay.com/WalletRestAPI/api/responses"
)

//ProceFundsTransfer - performs withdrawal operations
func ProceFundsTransfer(tx *gorm.DB, requestCreated *models.TransactionRequests) (responses.TransactionResponse, error) {
	//first get Balance of the DebitAccount
	response := responses.TransactionResponse{}
	cbalance1 := models.Accounts{}
	err := tx.Debug().Model(&models.Accounts{}).Where("account_no = ?", requestCreated.DebitAccount).Take(&cbalance1).Error
	if err != nil {
		return responses.TransactionResponse{}, err
	}
	//check if available balnce is enough for cash Withdrawal - should be more than the Amount and the Charge
	customerBal, _ := strconv.ParseFloat(cbalance1.AvailableBal, 64)
	withdrawalAmt, _ := strconv.ParseFloat(requestCreated.Amount, 64)
	if withdrawalAmt > customerBal {
		response.Procode = requestCreated.Procode
		response.ResponseCode = InsufficientBalance
		response.Remarks = InsufficientBalRemark
		response.Reference = requestCreated.TxnRef
		response.Amount = withdrawalAmt
		response.Account = cbalance1.AccountNo
		response.AvailableBalance = customerBal
		return response, nil
	}
	//do the withdrawal Debit entry - DR
	newBal := customerBal - withdrawalAmt
	//insert Debit entry
	t := models.Transactions{} //Transactions object
	t.MsgType = requestCreated.MsgType
	t.Procode = requestCreated.Procode
	t.Channel = requestCreated.Channel
	t.TxnRef = requestCreated.TxnRef
	t.AccountNo = requestCreated.DebitAccount
	t.Amount = requestCreated.Amount
	t.DrCr = "DR"
	t.Narration = requestCreated.Narration
	t.AvailableBal = fmt.Sprintf("%.2f", newBal)

	t.Prepare()
	t.Validate()
	_, err = t.SaveTransactions(tx)
	if err != nil {
		return responses.TransactionResponse{}, err
	}

	//do credit entry
	cbalance2 := models.Accounts{}
	//get Balance for the credit account_no
	err = tx.Debug().Model(&models.Accounts{}).Where("account_no = ?", requestCreated.CreditAccount).Take(&cbalance2).Error
	if err != nil {
		return responses.TransactionResponse{}, err
	}
	//Insert Credit entry
	receiverBal, _ := strconv.ParseFloat(cbalance2.AvailableBal, 64)
	receiverNewBal := receiverBal + withdrawalAmt

	t2 := models.Transactions{} //Transactions object
	t2.MsgType = requestCreated.MsgType
	t2.Procode = requestCreated.Procode
	t2.Channel = requestCreated.Channel
	t2.TxnRef = requestCreated.TxnRef
	t2.AccountNo = requestCreated.CreditAccount
	t2.Amount = requestCreated.Amount
	t2.DrCr = "CR"
	t2.Narration = requestCreated.Narration
	t2.AvailableBal = fmt.Sprintf("%.2f", receiverNewBal)

	t2.Prepare()
	t2.Validate()
	_, err = t2.SaveTransactions(tx)
	if err != nil {
		return responses.TransactionResponse{}, err
	}
	//Update balances on both credit and debit Accounts
	//1. debit account.
	updatedAcc1, err := UpdateBalance(tx, requestCreated.DebitAccount, newBal)
	if err != nil {
		return responses.TransactionResponse{}, err
	}
	//2.Credit Account
	updatedAcc2, err := UpdateBalance(tx, requestCreated.CreditAccount, receiverNewBal)
	if err != nil {
		return responses.TransactionResponse{}, err
	}

	response.Procode = requestCreated.Procode
	response.ResponseCode = Successful

	response.Reference = requestCreated.TxnRef
	response.Amount = withdrawalAmt

	if requestCreated.Procode == "210000" {
		response.AvailableBalance = receiverNewBal
		response.Account = updatedAcc2.AccountNo
		response.Remarks = "Cash Deposit Successful"
	} else {
		response.Remarks = "Cash Withdrawal Successful"
		response.AvailableBalance = newBal
		response.Account = updatedAcc1.AccountNo
	}

	return response, nil

}

//UpdateBalance Performs a balance update in Account Table
func UpdateBalance(tx *gorm.DB, account string, newBal float64) (*models.Accounts, error) {
	cbalance := models.Accounts{}
	err := tx.Debug().Model(&models.Accounts{}).Where("account_no = ?", account).Take(&cbalance).Error
	if err != nil {
		return &models.Accounts{}, err
	}
	cbalance.AvailableBal = fmt.Sprintf("%.2f", newBal)
	cbalance.ActualBal = fmt.Sprintf("%.2f", newBal)
	newCustomerAcc, err := cbalance.UpdateAccount(tx, account)
	if err != nil {
		return &models.Accounts{}, err
	}

	return newCustomerAcc, nil
}
