package modeltests

import (
	"log"
	"reflect"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/assert.v1"
	"tospay.com/WalletRestAPI/api/models"
)

func TestCustomer_SaveCustomer(t *testing.T) {

	// err := refreshCustomerAndAccountTable()
	// if err != nil {
	// 	log.Fatalf("Error refreshing  account and customer table %v\n", err)
	// }

	_, err := loadOneCustomer()
	if err != nil {
		log.Fatalf("Cannot laod customer %v\n", err)
	}

	newCustomer := models.Customer{
		ID:        2,
		FirstName: "Simon",
		LastName:  "Maingi",
		Email:     "simon@gmail.com",
		AccountNo: "254708003472",
	}
	savedCustomer, err := newCustomer.SaveCustomer(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the customer: %v\n", err)
		return
	}
	assert.Equal(t, newCustomer.ID, savedCustomer.ID)
	assert.Equal(t, newCustomer.FirstName, savedCustomer.FirstName)
	assert.Equal(t, newCustomer.LastName, savedCustomer.LastName)
	assert.Equal(t, newCustomer.AccountNo, savedCustomer.AccountNo)

}

func TestCustomer_FindAllCustomers(t *testing.T) {

	// err := refreshCustomerAndAccountTable()
	// if err != nil {
	// 	log.Fatalf("Error refreshing account and customer table %v\n", err)
	// }
	customer, err := loadOneCustomer()
	if err != nil {
		log.Fatalf("Error loading customer and account  table %v\n", err)
	}
	customers, err := customerInstance.FindAllCustomers(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the customers: %v\n", err)
		return
	}
	for _, v := range *customers {
		if !reflect.DeepEqual(v.AccountNo, customer.AccountNo) {
			t.Errorf("Customer.FindAllCustomers() = %v, want %v", customers, customer)
		}
	}

	//assert.Equal(t, len(*posts), 2)

}

func TestCustomer_UpdateCustomer(t *testing.T) {

	// err := refreshCustomerAndAccountTable()
	// if err != nil {
	// 	log.Fatalf("Error refreshing account and customer table: %v\n", err)
	// }
	customer, err := loadOneCustomer()
	if err != nil {
		log.Fatalf("Error loading table")
	}
	customerUpdate := models.Customer{
		ID:        1,
		FirstName: "Kamau",
		LastName:  "John",
		Email:     "John@gmail.com",
		AccountNo: "254712377789",
	}
	updatedCustomer, err := customerUpdate.UpdateCustomer(server.DB, customerUpdate.ID, customer.AccountNo)
	if err != nil {
		t.Errorf("this is the error updating the customer: %v\n", err)
		return
	}
	assert.Equal(t, updatedCustomer.ID, customerUpdate.ID)
	assert.Equal(t, updatedCustomer.FirstName, customerUpdate.FirstName)
	assert.Equal(t, updatedCustomer.LastName, customerUpdate.LastName)
	assert.Equal(t, updatedCustomer.AccountNo, customerUpdate.AccountNo)

}
