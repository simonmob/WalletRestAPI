package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"tospay.com/WalletRestAPI/api/authorization"
	"tospay.com/WalletRestAPI/api/models"
	"tospay.com/WalletRestAPI/api/responses"
)

//CreateCustomer creates a new Customer
func (server *Server) CreateCustomer(c *gin.Context) {

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

	customer := models.Customer{}
	err = c.BindJSON(&customer)
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
	}

	customer.Prepare()
	err = customer.Validate()
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
		return
	}

	customerCreated, err := customer.SaveCustomer(server.DB)

	if err != nil {

		//formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	//url := location.Get(c)
	//c.Writer.Header().Set("Location", fmt.Sprintf("%s%s/%s", url.Host, url.Path, customerCreated.AccountNo))
	//c.Header().Set("Location", fmt.Sprintf("%s%s/%s", c.Host, c.RequestURI, customerCreated.AccountNo))
	responses.JSON(c, http.StatusCreated, customerCreated)
}

//GetCustomers gets list of all registered customers
func (server *Server) GetCustomers(c *gin.Context) {

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
	customer := models.Customer{}

	customers, err := customer.FindAllCustomers(server.DB)
	if err != nil {
		responses.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(c, http.StatusOK, customers)
}

//GetCustomer gets customerdetails given the account no.
func (server *Server) GetCustomer(c *gin.Context) {

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

	account := c.Param("accountno")
	if account == "" {
		responses.ERROR(c, http.StatusBadRequest, errors.New("Account number required"))
		return
	}
	customer := models.Customer{}

	customerReceived, err := customer.FindCustomerByAccount(server.DB, account)
	if err != nil {
		responses.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(c, http.StatusOK, customerReceived)
}

//UpdateCustomer updates customer Details given the ID.
func (server *Server) UpdateCustomer(c *gin.Context) {

	id := c.Param("id")

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

	// Check if the customer exist
	customer := models.Customer{}
	err = server.DB.Debug().Model(models.Customer{}).Where("id = ?", id).Take(&customer).Error
	if err != nil {
		responses.ERROR(c, http.StatusNotFound, errors.New("Customer not found"))
		return
	}

	// Read the data posted and Start processing the request data
	customerUpdate := models.Customer{}
	err = c.BindJSON(&customerUpdate)
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
		return
	}

	customerUpdate.Prepare()
	err = customerUpdate.Validate()
	if err != nil {
		responses.ERROR(c, http.StatusUnprocessableEntity, err)
		return
	}

	customerUpdate.ID = customer.ID

	customerUpdated, err := customerUpdate.UpdateCustomer(server.DB, customerUpdate.ID, customer.AccountNo)

	if err != nil {
		//formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(c, http.StatusOK, customerUpdated)
}
