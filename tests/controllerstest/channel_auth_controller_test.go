package controllerstest

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/assert.v1"
)

func TestServer_ChannelAuth(t *testing.T) {

	_, err := loadOneChannel()
	if err != nil {
		log.Fatalf("Cannot load channel %v\n", err)
	}

	tests := []struct {
		name         string
		inputJSON    string
		statusCode   int
		errorMessage string
	}{
		{
			name:         "ChannelAuthSuccessTest",
			inputJSON:    `{"channel": "USSD", "description": "test channel"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			name:         "ChannelAuthFailTest",
			inputJSON:    `{"channel": "", "description": "test channel"}`,
			statusCode:   422,
			errorMessage: "Channel Required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body := bytes.NewBufferString(tt.inputJSON)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest("POST", "/channelAuth", body)
			c.Request.Header.Set("Content-Type", "application/json")
			//c.Request.Header.Set("Authorization", tt.tokenGiven)

			server.ChannelAuth(c)

			assert.Equal(t, w.Code, tt.statusCode)
			if tt.statusCode == 200 {
				assert.NotEqual(t, w.Body.String(), "")
			}
			if tt.statusCode == 422 && tt.errorMessage != "" {
				responseMap := make(map[string]interface{})
				err := json.Unmarshal([]byte(w.Body.String()), &responseMap)
				if err != nil {
					t.Errorf("Cannot convert to json: %v", err)
				}
				if !reflect.DeepEqual(responseMap["error"], tt.errorMessage) {
					t.Errorf("ChannelAuth response= %v, want %v", responseMap["error"], tt.errorMessage)
				}
				//assert.Equal(t, responseMap["error"], tt.errorMessage)
			}
		})
	}
}

func TestServer_GetChannelToken(t *testing.T) {

	_, err := loadOneChannel()
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		name       string //given name for the test.
		channelstr string
		want       string
		wantErr    bool
	}{
		{
			name:       "GetchannelTokenSuccessfulTest",
			channelstr: "USSD",
			want:       "token",
			wantErr:    false,
		},
		{
			name:       "GetchannelTokenFailTest",
			channelstr: "",
			want:       "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := server.GetChannelToken(tt.channelstr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.GetChannelToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				if len(got) > 20 {

				} else {
					t.Errorf("Server.GetChannelToken() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
