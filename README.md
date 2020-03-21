## Demo Wallet Rest API
It uses
1. Gorm -  ORM library for Golang
2. go modules - go packages dependency management
3. postgres
4. jwt - for channel authentication
5. Gin HTTP web framework

## Functionalities
1. Channel creation
2. Customer creation
3. Balance enquiry
4. Mini-statement Enquiry
5. Cash Deposit
6. Cash withdrawal

## API Security
It uses jwt for authorization. Note: the authorization is done for channels. only authorized channels can post, put or get. The token is valid for 24 hours.

## How to configure and Test the API
1. Create a db named walletapi in Postgres
2. Set the DB Configs in the `.env` file. e.g
```ENV
DB_IP=127.0.0.1
DB_DRIVER=postgres
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=walletapi
DB_PORT=5439 #User defined - The Default postgres port is 5432
```
3. Set the CashAccount and api port values in the config
```ENV
#cash account
CASH_ACCOUNT=254712345678
#api port
API_PORT=8083
```
3. Run the Api. Navigate to WalletRestAPI folder and start using go run command.
```command
 go run main.go
 ```
4. Database objects will be auto migrated to the postgres DB
```list
1 accounts
2 channel
3 customer
4 transaction_requests
5 transactions
```
6. Default channel (USSD) and Account(251712345678- for cash account) will be inserted automatically after models AutoMigrate

7. Test the api end points

## API endpoints
###### 1. POST -[/channels](localhost:8083/channels) - this end point creates new channel.
* Request
```JSON
{
	"channel": "USSD",
	"description":"USSD channel"
}
```
* Response
```JSON
{
    "id": 6,
    "channel": "MOBILE",
    "description": "MOBILE channel",
    "created_at": "2020-03-15T17:47:21.0290587+03:00",
    "updated_at": "2020-03-15T17:47:21.0290587+03:00"
}
```

###### 2. GET -[/GetChannels]() - Gets a list of the created channels.
* Request - [/GetChannels]()
```JSON
{
       "id": 1,
       "channel": "USSD",
       "description": "ussd channel",
       "created_at": "2020-03-15T17:01:56.522541+03:00",
       "updated_at": "2020-03-15T17:01:56.522541+03:00"
   },
   {
       "id": 4,
       "channel": "WEB",
       "description": "WEB channel",
       "created_at": "2020-03-15T17:06:28.39293+03:00",
       "updated_at": "2020-03-15T17:06:28.39293+03:00"
   },
   {
       "id": 6,
       "channel": "MOBILE",
       "description": "MOBILE channel",
       "created_at": "2020-03-15T17:47:21.029059+03:00",
       "updated_at": "2020-03-15T17:47:21.029059+03:00"
   }
```

###### 3. POST - [/channelAuth]() - Gets a token for authorization. This token should be passed on the header of the all the other requests.
* Requests
```JSON
{
	"channel": "MOBILE"
}
```
* Response
```JSON
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE1ODQzNzA0NDQsInVzZXJfaWQiOiJNT0JJTEUifQ.4VOl6gzybdGLOVYJAbTs4kF-2s8DJltLM4qlDNAsB3s"
}
```

###### 4. POST - [/customers/create]() - creates a new customer. A call to this end point opens an account for the customer. Note the account should be a Phone Number format starting with 254xxxxxxxxx. New accounts in accounts table will have a balance of 0.00.
* Request
```JSON
{
	"firstname":"Simon",
	"lastname":"Maingi",
	"email":"kmaing.simon@gmail.com",
	"account_no":"254708003472"
}
```
* Response
```JSON
{
    "id": 1,
    "firstname": "Simon",
    "lastname": "Maingi",
    "email": "kmaing.simon@gmail.com",
    "account_no": "254708003472",
    "created_at": "2020-03-15T18:37:32.683894+03:00",
    "updated_at": "2020-03-15T18:37:32.683894+03:00"
}
```

###### 5. GET - [/customers/get]() - gets a list of all customers. Limit of 100.
* Response
```JSON
[
    {
        "id": 1,
        "firstname": "Simon",
        "lastname": "Maingi",
        "email": "kmaing.simon@gmail.com",
        "account_no": "254708003472",
        "created_at": "2020-03-15T18:53:58.069803+03:00",
        "updated_at": "2020-03-15T18:53:58.069803+03:00"
    },
    {
        "id": 2,
        "firstname": "John",
        "lastname": "Mwangi",
        "email": "mwangi@gmail.com",
        "account_no": "254712345678",
        "created_at": "2020-03-15T18:55:38.742371+03:00",
        "updated_at": "2020-03-15T18:55:38.742371+03:00"
    }
]
```

###### 6. GET-[/customers/get/{accountno}]() - Get customer details of the given account number. No request Body
* Get - [/customers/get/254708003472]()
* Resposne
```JSON
{
    "id": 1,
    "firstname": "Simon",
    "lastname": "Maingi",
    "email": "kmaing.simon@gmail.com",
    "account_no": "254708003472",
    "created_at": "2020-03-15T18:53:58.069803+03:00",
    "updated_at": "2020-03-15T18:53:58.069803+03:00"
}
```

###### 7. PUT - [/customers/update/{id}]() - Updates account details of the given Customer id. Id can be found after get details above
* Request -[/customers/update/2]()
```JSON
{
	"firstname":"John",
	"lastname":"Kamau",
	"email":"kamau.john@gmail.com",
	"account_no":"254712345678"
}
```
* Resposne
```JSON
{
    "id": 2,
    "firstname": "John",
    "lastname": "Kamau",
    "email": "kamau.john@gmail.com",
    "account_no": "254712345678",
    "created_at": "2020-03-15T18:55:38.742371+03:00",
    "updated_at": "2020-03-15T19:28:55.931394+03:00"
}
```

###### 8. POST - [/transactions]() - processes all transaction related requests. The transaction type is identified by the value of the supplied processing code in the request.
a. Balance Enquiry - processing code `310000`. debit account should be the account you want to request balance for.

* Request
```JSON
{
  "msg_type":"0200",
  "pro_code":"310000",
  "channel":"USSD",
  "txn_ref":"79827162829",
  "amount":"0.00",
  "narration":"Balance Enquiry",
  "debit_account":"254708003472",
  "credit_account": ""
}
```
* Response
```JSON
{
  "procode": "310000",
  "response_code": "000",
  "remarks": "Balance Enquiry Successful",
  "reference": "79827162829",
  "amount": 0,
  "account_no": "254708003472",
  "available_bal": 0
}
```

b. Cash Deposit - processing code `210000`. The cash account(254712345678) is used as the Debit account in this case.

* Request
```JSON
{
  "msg_type":"0200",
  "pro_code":"210000",
  "channel":"USSD",
  "txn_ref":"79827162824",
  "amount":"1000.50",
  "narration":"Cash Deposit",
  "debit_account":"254708003472",
  "credit_account": "254708003472"
}      
```
* Response
```JSON
{
  "procode": "210000",
  "response_code": "000",
  "remarks": "Cash Deposit Successful",
  "reference": "79827162824",
  "amount": 1000.5,
  "account_no": "254708003472",
  "available_bal": 1000.5
}
```

c. Cash withdrawal - processing code `010000`. Credit account here is the cash account (`254712345678`).

* Request
```JSON
{
	"msg_type":"0200",
	"pro_code":"010000",
	"channel":"USSD",
	"txn_ref":"79827162825",
	"amount":"200",
	"narration":"Cash Withdrawal",
	"debit_account":"254708003472",
	"credit_account": ""
}
```
* Response
```JSON
{
  "procode": "010000",
  "response_code": "000",
  "remarks": "Cash Withdrawal Successful",
  "reference": "79827162825",
  "amount": 200,
  "account_no": "254708003472",
  "available_bal": 800.5
}
```

d. Mini-statement Enquiry - processing code `380000`. Amount should be zero and debit account should be the one you want to request a ministatement for.

* Request
```JSON
{
	"msg_type":"0200",
	"pro_code":"380000",
	"channel":"USSD",
	"txn_ref":"79827162826",
	"amount":"0.00",
	"narration":"Ministatement Enquiry",
	"debit_account":"254708003472",
	"credit_account": ""
}
```
* Response
```JSON
{
  "procode": "380000",
  "response_code": "000",
  "remarks": "Ministatement Enquiry Successful",
  "reference": "79827162826",
  "amount": 0,
  "account_no": "",
  "available_bal": 0,
  "Ministatement": [
	  {
		  "Txn_Date_time": "2020-03-16T00:08:06.430269+03:00",
		  "Txn_Type": "Cash Withdrawal",
		  "amount": "200.00",
		  "dr_cr": "DR"
	  },
	  {
		  "Txn_Date_time": "2020-03-16T00:07:44.370908+03:00",
		  "Txn_Type": "Cash Deposit",
		  "amount": "1000.50",
		  "dr_cr": "CR"
	  }
  ]
}
```


## Unit Tests

#### 1. Models Tests - found in modeltests folder under tests

To run all the tests at once use; `go test -v` on the modeltests folder.

a. SaveChannel test - This runs a successful Channel Save test. Found in `channel_test.go` file. To test, run `go test -v --run TestChannel_SaveChannel` command.

Test Data;
```js
ID:          1,
Channel:     "USSD",
Description: "test Channel",
CreatedAt:   time.Now(),
UpdatedAt:   time.Now(),
```
Assertions
- Check if the returned `Channel` name is the same as the given test channel Name.
```js
if !reflect.DeepEqual(savedChannel.Channel, newChannel.Channel) {
  t.Errorf("Channel.SaveChannel() = %v, want %v", savedChannel, newChannel)
}
```
- Check if the returned `Description` is the same as the given test Description.
```js
if !reflect.DeepEqual(savedChannel.Description, newChannel.Description) {
  t.Errorf("Channel.SaveChannel() = %v, want %v", savedChannel, newChannel)
}
```

b. FindAllChannels -Gets a list of saved channels. Found in `channel_test.go` file. To test, run `go test -v --run TestChannel_FindAllChannels` command.

Test Data.

`loadOneChannel()` - saves one channel.

Assertion
- Check if the returned list has a `Channel` value same as the one saved in the test data.
```js
if !reflect.DeepEqual("USSD", v.Channel) {
  t.Errorf("Channel.FindAllChannels() got = %v, want %v", v.Channel, "USSD")
}
```

c. SaveCustomer - Saves a new customer. Found in `customer_test.go` file. To test, run `go test -v --run TestCustomer_SaveCustomer` command.

Test Data.
```js
ID:        2,
FirstName: "Simon",
LastName:  "Maingi",
Email:     "simon@gmail.com",
AccountNo: "254708003472",
```
Assertions
- Check if returned `ID,FirstName,LastName and AccountNo` are the same as the given test Data
```js
assert.Equal(t, newCustomer.ID, savedCustomer.ID)
assert.Equal(t, newCustomer.FirstName, savedCustomer.FirstName)
assert.Equal(t, newCustomer.LastName, savedCustomer.LastName)
assert.Equal(t, newCustomer.AccountNo, savedCustomer.AccountNo)
```

d. FindAllCustomers - Gets a list of the saved customers. Found in `customer_test.go` file. To test, run `go test -v --run TestCustomer_FindAllCustomers` command.

Test Data.

`loadOneCustomer` - saves one customer.

Assertions
- Check if the returned `AccountNo` is the same given in the test data.
```js
if !reflect.DeepEqual(v.AccountNo, customer.AccountNo) {
  t.Errorf("Customer.FindAllCustomers() = %v, want %v", customers, customer)
}
```

e. UpdateCustomer - Updates customer details. Found in `customer_test.go` file. To test, run `go test -v --run TestCustomer_FindAllCustomers` command.

Test Data.

`loadOneCustomer` - saves one customer. This is the customer to be updated with the details below.
```js
ID:        1,
FirstName: "Kamau",
LastName:  "John",
Email:     "John@gmail.com",
AccountNo: "254712377789",
```
Assertions
- Check if returned updated `ID,FirstName,LastName and AccountNo` are the same as the given update test Data.
```js
assert.Equal(t, updatedCustomer.ID, customerUpdate.ID)
assert.Equal(t, updatedCustomer.FirstName, customerUpdate.FirstName)
assert.Equal(t, updatedCustomer.LastName, customerUpdate.LastName)
assert.Equal(t, updatedCustomer.AccountNo, customerUpdate.AccountNo)
```

#### 2. Controllers Tests - found in controllerstest folder under tests.
To run all the tests at once use; `go test -v` on the controllerstest folder.

a. ChannelAuth - Checks if a channel is authorized to access the API. Found in `channel_auth_controller_test.go` file. To test, run `go test -v --run TestServer_ChannelAuth` command.

Test Data.

```js
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
```
Assertions
- Check statusCode returned, if `200`, the channel is authorized. This is as expected in test data one above.
```js
assert.Equal(t, w.Code, tt.statusCode)
if tt.statusCode == 200 {
  assert.NotEqual(t, w.Body.String(), "")
}
```
- Check statusCode returned, if `422`, the channel is not authorized and an error message should be returned `Channel Required` as expected in test data two.
```js
if tt.statusCode == 422 && tt.errorMessage != "" {
  responseMap := make(map[string]interface{})
  err := json.Unmarshal([]byte(w.Body.String()), &responseMap)
  if err != nil {
    t.Errorf("Cannot convert to json: %v", err)
  }
  if !reflect.DeepEqual(responseMap["error"], tt.errorMessage) {
    t.Errorf("ChannelAuth response= %v, want %v", responseMap["error"], tt.errorMessage)
  }
  ```

b. CreateCustomer - Creates a customer.Found in `customers_controller_test.go` file. To test, run `go test -v --run TestServer_CreateCustomer` command.

Test Data.
```js
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
```
Assertions
- Check for customer creation successful test. If statusCode code returned is `201`, means customer was created successfully and error message is nil as expected in the test data one. Aslo checks if the returned customer `AccountNo` is same as the one given in the test data.
```js
assert.Equal(t, w.Code, tt.statusCode)
if tt.statusCode == 201 {
  assert.Equal(t, responseMap["account_no"], "254708003472")
}
```
- Check if channel is authorized to do customer creation. If stat
 returned is `401` and the erro message is `Unauthorized` as expected in the given test data. This is asserted by supplying wrong token in the request.
 ```js
 if tt.statusCode == 401 && tt.errorMessage != "" {
   assert.Equal(t, responseMap["error"], tt.errorMessage)
 }
 ```


 c. Processes transactions IncomingRequest - Processes all incoming transaction requests. Found in `transaction_requests_controller_test.go` file. To test, run `go test -v --run TestServer_IncomingRequest` command.

 Test Data
 ```js
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
 ```
 Assertions
 - Check successful balance Enquiry,Successful cash Deposit,successful cash withdrawal and successful Ministatement enquiry. Checks if the returned response_code code is `000` as expected in the given test datas. The http statusCode is `200`.
 ```js
 assert.Equal(t, w.Code, tt.statusCode)
 if tt.statusCode == 200 && responseMap["response_code"] == "000" {
   assert.Equal(t, responseMap["response_code"], "000")

 }
 ```
 - Check for Unauthorized transaction. Status is `401` and errorMessage `Unauthorized` as expected in the test data.
 ```js
 if tt.statusCode == 401 && tt.errorMessage != "" {
   assert.Equal(t, responseMap["error"], tt.errorMessage)
 }
 ```
 - Check for InsufficientBalance. ResponseCode `020` as expected in the test data.
 ```js
 else if responseMap["response_code"] == "020" {
   assert.Equal(t, responseMap["response_code"], "020")
 }
 ```
