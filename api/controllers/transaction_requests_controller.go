package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"tospay.com/WalletRestAPI/api/authorization"
	"tospay.com/WalletRestAPI/api/models"
	"tospay.com/WalletRestAPI/api/responses"
	"tospay.com/WalletRestAPI/api/utils"
)

//IncomingRequest entry point for all transaction Requests
func (server *Server) IncomingRequest(c *gin.Context) {
	//CHeck if the auth token is valid and  get the channel from it
	channel, err := authorization.ExtractTokenID(c)
	if err != nil {
		responses.ERROR(c, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if channel == "" {
		responses.ERROR(c, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	request := models.TransactionRequests{}
	err = c.BindJSON(&request) //get request values
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
	}

	request.Prepare()
	err = request.Validate()
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
		return
	}

	procode := request.Procode
	//Set the Debit and credit accounts for Cash withdrawal and Cash Deposit
	switch procode {
	case "010000": //Cash withdrawal
		//1. DebitAccount = Customer Account
		//2. CreditAccount = CashAccount
		if request.CreditAccount == "" || request.CreditAccount != "254712345678" {
			request.CreditAccount = server.Configs.CashAccount
		}

		break
	case "210000": //cash Deposit
		//1. DebitAccount = CashAccount
		//2. CreditAccount = Customer Account
		if request.DebitAccount == "" || request.DebitAccount != "254712345678" {
			request.DebitAccount = server.Configs.CashAccount
		}
	}
	//log the message request in DB. this step works like messages logger.
	requestCreated, err := request.SaveTransactionRequest(server.DB)
	if err != nil {
		responses.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	//start DB transaction to process incoming requests
	//This is important to rollback changes which were recorded in DB if an error occurs.
	tx := server.DB.Begin()
	defer func() { //rolls back what has happened and update response if error occurs during the txn life cycle
		if r := recover(); r != nil {
			tx.Rollback()
			request.ResponseCode = utils.Failed
			request.Remarks = utils.ErrorRemark
			_, _ = request.UpdateTransactionResponse(server.DB, request.TxnRef)
		}
	}()
	//Process request depending on processing Code
	response := responses.TransactionResponse{}
	switch procode {
	case "010000": //Cash Withdrawal
		response, err = utils.ProceFundsTransfer(tx, requestCreated)
		if err != nil {
			tx.Rollback()
			request.ResponseCode = utils.Failed
			request.Remarks = utils.ErrorRemark
			_, _ = request.UpdateTransactionResponse(server.DB, request.TxnRef)
			responses.ERROR(c, http.StatusInternalServerError, err)
			return
		}
	case "210000": //Cash Deposit
		response, err = utils.ProceFundsTransfer(tx, requestCreated)
		if err != nil {
			tx.Rollback()
			request.ResponseCode = utils.Failed
			request.Remarks = utils.ErrorRemark
			_, _ = request.UpdateTransactionResponse(server.DB, request.TxnRef)
			responses.ERROR(c, http.StatusInternalServerError, err)
			return
		}

	case "310000": //Balance enquiry
		response, err = utils.GetBalance(tx, requestCreated)
		if err != nil {
			tx.Rollback()
			request.ResponseCode = utils.Failed
			request.Remarks = utils.ErrorRemark
			_, _ = request.UpdateTransactionResponse(server.DB, request.TxnRef)
			responses.ERROR(c, http.StatusInternalServerError, err)
			return
		}
	case "380000":
		response, err = utils.GetMinistatement(tx, requestCreated)
		if err != nil {
			tx.Rollback()
			request.ResponseCode = utils.Failed
			request.Remarks = utils.ErrorRemark
			_, _ = request.UpdateTransactionResponse(server.DB, request.TxnRef)
			responses.ERROR(c, http.StatusInternalServerError, err)
			return
		}
	default:
		request.ResponseCode = utils.Failed
		request.Remarks = utils.ProcodeNotDefined
		_, _ = request.UpdateTransactionResponse(server.DB, request.TxnRef)
		responses.ERROR(c, http.StatusInternalServerError, errors.New("Processing Code not Defined"))
		return
	}
	tx.Commit() //Commit the transaction if no error occured
	request.ResponseCode = response.ResponseCode
	request.Remarks = response.Remarks
	_, _ = request.UpdateTransactionResponse(server.DB, request.TxnRef)
	responses.JSON(c, http.StatusOK, response)
}
