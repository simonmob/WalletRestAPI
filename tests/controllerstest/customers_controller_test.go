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

func TestServer_CreateCustomer(t *testing.T) {
	err := refreshCustomerAndAccountTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = loadOneCustomer()
	if err != nil {
		log.Fatalf("Cannot load Customer %v\n", err)
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
		errorMessage string
	}{
		{
			name:         "CreateCustomer_Success",
			inputJSON:    `{"firstname":"Simon","lastname":"Maingi","email":"kmaing.simon@gmail.com","account_no":"254708003472"}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			name:         "CreateCustomer_Fail",
			inputJSON:    `{"firstname":"Simon","lastname":"Maingi","email":"kmaing.simon@gmail.com","account_no":"254708003472"}`,
			statusCode:   401,
			tokenGiven:   "wrong token",
			errorMessage: "Unauthorized",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body := bytes.NewBufferString(tt.inputJSON)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest("POST", "/customers/create", body)
			c.Request.Header.Set("Content-Type", "application/json")
			c.Request.Header.Set("Authorization", tt.tokenGiven)

			server.CreateCustomer(c)
			//fmt.Println("Response:" + w.Body.String())
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(w.Body.String()), &responseMap)
			if err != nil {
				fmt.Printf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, w.Code, tt.statusCode)
			if tt.statusCode == 201 {
				assert.Equal(t, responseMap["account_no"], "254708003472")
			}
			if tt.statusCode == 401 && tt.errorMessage != "" {
				assert.Equal(t, responseMap["error"], tt.errorMessage)
			}

		})
	}
}
