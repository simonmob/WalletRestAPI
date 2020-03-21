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
	err = c.BindJSON(&request)
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
	//Only insert is required here. this step works like messages logger.
	switch procode {
	case "010000":
		//1. DebitAccount = Customer Account
		//2. CreditAccount = CashAccount
		if request.CreditAccount == "" {
			request.CreditAccount = server.Configs.CashAccount
		}

		break
	case "210000":
		//1. DebitAccount = CashAccount
		//2. CreditAccount = Customer Account
		if request.DebitAccount == "" {
			request.DebitAccount = server.Configs.CashAccount
		}
	}
	requestCreated, err := request.SaveTransactionRequest(server.DB)
	if err != nil {
		responses.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	//start DB transaction to process incoming requests
	tx := server.DB.Begin()
	defer func() { //rolls back what has happened and update resposne if error occurs during the txn life cycle
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
	tx.Commit()
	request.ResponseCode = response.ResponseCode
	request.Remarks = response.Remarks
	_, _ = request.UpdateTransactionResponse(server.DB, request.TxnRef)
	responses.JSON(c, http.StatusOK, response)
}
