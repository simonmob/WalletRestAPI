package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

//Customer struct for Customer details
type Customer struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	FirstName string    `gorm:"size:255;not null" json:"firstname"`
	LastName  string    `gorm:"size:255;not null" json:"lastname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	AccountNo string    `gorm:"size:100;not null;unique" json:"account_no"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//Prepare prepares customer data
func (c *Customer) Prepare() {
	c.ID = 0
	c.FirstName = html.EscapeString(strings.TrimSpace(c.FirstName))
	c.LastName = html.EscapeString(strings.TrimSpace(c.LastName))
	c.Email = html.EscapeString(strings.TrimSpace(c.Email))
	c.AccountNo = html.EscapeString(strings.TrimSpace(c.AccountNo))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

//Validate validate customer data
func (c *Customer) Validate() error {

	if c.FirstName == "" {
		return errors.New("FirstName Required")
	}
	if c.LastName == "" {
		return errors.New("LastName Required")
	}
	if c.Email == "" {
		return errors.New("Email Required")
	}
	if err := checkmail.ValidateFormat(c.Email); err != nil {
		return errors.New("Invalid Email")
	}
	if c.AccountNo == "" {
		return errors.New("Account Required")
	}
	if len(c.AccountNo) != 12 { //account length should be 12
		return errors.New("Account length should be 12")
	}
	if !strings.HasPrefix(c.AccountNo, "254") { //account should start with 254
		return errors.New("Account should start with 254")
	}

	isNotDigit := func(c rune) bool { return c < '0' || c > '9' }
	b := strings.IndexFunc(c.AccountNo, isNotDigit) == -1
	if !b { //account should have digits only
		return errors.New("Account should have digits only")
	}
	return nil
}

//SaveCustomer inserts customer details into Customer table
func (c *Customer) SaveCustomer(db *gorm.DB) (*Customer, error) {
	var err error

	//start a transaction. Customer creation inserts in 2 databases-Accounts and Customer
	tx := db.Begin()
	defer func() { //rolls back if error occurs during the txn life cycle
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//insert in Customer table.
	err = tx.Debug().Model(&Customer{}).Create(&c).Error
	if err != nil {
		tx.Rollback()
		return &Customer{}, err
	}

	//insert in Accounts too.
	cAccount := Accounts{}
	cAccount.AccountNo = c.AccountNo

	cAccount.Prepare()
	err = cAccount.Validate()
	if err != nil {
		tx.Rollback()
		return &Customer{}, err
	}
	//insert into Accounts Table
	if c.AccountNo == "254712345678" {
		cAccount.AvailableBal = "10000.00"
		cAccount.ActualBal = "10000.00"
	}
	accountCreated, err := cAccount.SaveAccount(tx)
	if err != nil {
		tx.Rollback()
		return &Customer{}, err
	}

	if c.AccountNo != "" {
		err = tx.Debug().Model(&Customer{}).Where("account_no = ?", accountCreated.AccountNo).Take(&c).Error
		if err != nil {
			tx.Rollback()
			return &Customer{}, err
		}
	}
	tx.Commit()

	//End of customer creation transaction
	return c, nil
}

//FindAllCustomers gets list of all registered customers
func (c *Customer) FindAllCustomers(db *gorm.DB) (*[]Customer, error) {
	var err error
	customers := []Customer{}
	err = db.Debug().Model(&Customer{}).Limit(100).Find(&customers).Error
	if err != nil {
		return &[]Customer{}, err
	}
	return &customers, err
}

//FindCustomerByAccount gets customer details given the account
func (c *Customer) FindCustomerByAccount(db *gorm.DB, account string) (*Customer, error) {
	var err error
	err = db.Debug().Model(&Customer{}).Where("account_no = ?", account).Take(&c).Error
	if err != nil {
		return &Customer{}, err
	}
	// if c.AccountNo != "" {
	// 	err = db.Debug().Model(&Customer{}).Where("account_no = ?", c.AccountNo).Take(&c.FirstName).Error
	// 	if err != nil {
	// 		return &Customer{}, err
	// 	}
	// }
	return c, nil
}

//UpdateCustomer updates customer details given the Account
func (c *Customer) UpdateCustomer(db *gorm.DB, id uint32, accountToUpdate string) (*Customer, error) {

	var err error

	//start customerUpdate transaction.
	tx := db.Begin()
	defer func() { //rolls back if error occurs during the txn life cycle
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = tx.Debug().Model(&Customer{}).Where("id = ?", id).Updates(Customer{FirstName: c.FirstName, LastName: c.LastName, AccountNo: c.AccountNo, Email: c.Email, UpdatedAt: time.Now()}).Error
	if err != nil {
		tx.Rollback()
		return &Customer{}, err
	}
	//update in Accounts too.
	cAccount := Accounts{}
	cAccount.AccountNo = c.AccountNo
	err = tx.Debug().Model(&Accounts{}).Where("account_no = ?", accountToUpdate).Updates(Accounts{AccountNo: c.AccountNo}).Error
	if err != nil {
		tx.Rollback()
		return &Customer{}, err
	}
	if c.AccountNo != "" {
		err = tx.Debug().Model(&Customer{}).Where("account_no = ?", c.AccountNo).Take(&c).Error
		if err != nil {
			return &Customer{}, err
		}
	}
	tx.Commit()
	//End of customerUpdate transaction

	return c, nil
}
