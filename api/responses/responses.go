package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//JSON formats responses to json
func JSON(c *gin.Context, statusCode int, data interface{}) {
	// w.WriteHeader(statusCode)
	// err := json.NewEncoder(w).Encode(data)
	// if err != nil {
	// 	fmt.Fprintf(w, "%s", err.Error())
	// }

	//c.Render(statusCode, render.IndentedJSON{Data: data})
	c.JSON(statusCode, data)
}

//ERROR formats error responses to json
func ERROR(c *gin.Context, statusCode int, err error) {
	if err != nil {
		JSON(c, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(c, http.StatusBadRequest, nil)
}

//TransactionResponse struct hold transaction responses
type TransactionResponse struct {
	Procode          string          `json:"procode"`
	ResponseCode     string          `json:"response_code"`
	Remarks          string          `json:"remarks"`
	Reference        string          `json:"reference"`
	Amount           float64         `json:"amount"`
	Account          string          `json:"account_no"`
	AvailableBalance float64         `json:"available_bal"`
	Ministatement    []Ministatement `json:",omitempty,ministatement"`
}

//Ministatement struct
type Ministatement struct {
	Created_At string `json:"Txn_Date_time"`
	Narration  string `json:"Txn_Type"`
	Amount     string `json:"amount"`
	DR_CR      string `json:"dr_cr"`
}

//TxnResponse format transaction response to json
func TxnResponse(c *gin.Context, statusCode int, tr TransactionResponse) {
	JSON(c, statusCode, tr)
}
