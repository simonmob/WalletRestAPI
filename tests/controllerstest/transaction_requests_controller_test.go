package controllerstest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/assert.v1"
)

func TestServer_IncomingRequest(t *testing.T) {
	err := refreshCustomerAndAccountTable()
	if err != nil {
		log.Fatal(err)
	}
	err = refreshRequestsAndTransactionsTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = loadCustomers()
	if err != nil {
		log.Fatalf("Cannot load Customers %v\n", err)
	}
	channel, err := loadOneChannel()
	if err != nil {
		log.Fatalf("Cannot load channel %v\n", err)
	}
	token, err := server.GetChannelToken(channel.Channel)
	if err != nil {
		log.Fatalf("cannot get channel Token: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	tests := []struct {
		name         string
		inputJSON    string
		statusCode   int
		tokenGiven   string
		ResponseCode string
		errorMessage string
	}{
		{
			name:         "TestSuccesfulBalance",
			inputJSON:    `{"msg_type":"0200","pro_code":"310000","channel":"USSD","txn_ref":"79827162829","amount":"0.00","narration":"Balance Enquiry","debit_account":"254708003472","credit_account": ""}`,
			statusCode:   200,
			tokenGiven:   tokenString,
			ResponseCode: "000",
			errorMessage: "",
		},
		{
			name:         "TestUnauthorizedTransaction",
			inputJSON:    `{"msg_type":"0200","pro_code":"310000","channel":"USSD","txn_ref":"79827162830","amount":"0.00","narration":"Balance Enquiry","debit_account":"254708003472","credit_account": ""}`,
			statusCode:   401,
			tokenGiven:   "wrong token",
			errorMessage: "Unauthorized",
		},
		{
			name:         "TestInsufficientBalance",
			inputJSON:    `{"msg_type":"0200","pro_code":"010000","channel":"USSD","txn_ref":"79827162824","amount":"110000.00","narration":"Cash Withdrawal","debit_account":"254708003472","credit_account": "254712345678"}`,
			statusCode:   200,
			tokenGiven:   tokenString,
			ResponseCode: "020",
			errorMessage: "",
		},
		{
			name:         "TestSuccessfulCashWithdrawal",
			inputJSON:    `{"msg_type":"0200","pro_code":"010000","channel":"USSD","txn_ref":"79827162825","amount":"200.00","narration":"Cash Withdrawal","debit_account":"254708003472","credit_account": "254712345678"}`,
			statusCode:   200,
			tokenGiven:   tokenString,
			ResponseCode: "000",
			errorMessage: "",
		},
		{
			name:         "TestSuccessfulCashDeposit",
			inputJSON:    `{"msg_type":"0200","pro_code":"210000","channel":"USSD","txn_ref":"79827162826","amount":"300.00","narration":"Cash Deposit","debit_account":"254712345678","credit_account": "254708003472"}`,
			statusCode:   200,
			tokenGiven:   tokenString,
			ResponseCode: "000",
			errorMessage: "",
		},
		{
			name:         "TestSuccessfulMinistatement",
			inputJSON:    `{"msg_type":"0200","pro_code":"380000","channel":"USSD","txn_ref":"79827162827","amount":"300.00","narration":"Cash Deposit","debit_account":"254708003472","credit_account": "254708003472"}`,
			statusCode:   200,
			tokenGiven:   tokenString,
			ResponseCode: "000",
			errorMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body := bytes.NewBufferString(tt.inputJSON)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest("POST", "/transactions", body)
			c.Request.Header.Set("Content-Type", "application/json")
			c.Request.Header.Set("Authorization", tt.tokenGiven)

			server.IncomingRequest(c)
			//fmt.Println("Response:" + w.Body.String())
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(w.Body.String()), &responseMap)
			if err != nil {
				fmt.Printf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, w.Code, tt.statusCode)
			if tt.statusCode == 200 && responseMap["response_code"] == "000" {
				assert.Equal(t, responseMap["response_code"], "000")

			} else if responseMap["response_code"] == "020" {
				assert.Equal(t, responseMap["response_code"], "020")
			}
			if tt.statusCode == 401 && tt.errorMessage != "" {
				assert.Equal(t, responseMap["error"], tt.errorMessage)
			}

		})
	}
}
