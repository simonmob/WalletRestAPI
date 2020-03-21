package sampledata

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"tospay.com/WalletRestAPI/api/models"
)

var channel = models.Channel{
	Channel:     "USSD",
	Description: "Test USSD Channel",
}

var cashAccountCustomer = models.Customer{
	FirstName: "Cash",
	LastName:  "Account",
	Email:     "cash@gmail.com",
	AccountNo: "254712345678",
}

//Load Inserts test channel and cash account.
func Load(db *gorm.DB) {

	//drop if exists
	err := db.Debug().DropTableIfExists(&models.Channel{}, &models.Customer{}, &models.Accounts{}, &models.TransactionRequests{}, &models.Transactions{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Channel{}, &models.Accounts{}, &models.Customer{}, &models.TransactionRequests{}, &models.Transactions{}).Error //database migration
	//err = db.Debug().AutoMigrate(&models.Channel{}, &models.Customer{}, &models.Accounts{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	//insert Channels
	channel.Prepare()
	err = channel.Validate("")
	if err != nil {
		fmt.Println(err)
		return
	}

	channelCreated, err := channel.SaveChannel(db)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Channel Created:")
	fmt.Println(channelCreated)

	//Insert cash account
	cashAccountCustomer.Prepare()
	err = cashAccountCustomer.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	customerCreated, err := cashAccountCustomer.SaveCustomer(db)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Cash Account Created:")
	fmt.Println(customerCreated)
}
