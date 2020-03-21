package controllerstest

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"tospay.com/WalletRestAPI/api/controllers"
	"tospay.com/WalletRestAPI/api/models"
)

var server = controllers.Server{}

var channelInstance = models.Channel{}
var customerInstance = models.Customer{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	InitializeDB()

	os.Exit(m.Run())

	//sampledata.Load(server.DB)

	//	server.Run(":" + os.Getenv("API_PORT"))
}

func InitializeDB() {
	var err error
	Dbdriver := os.Getenv("DB_DRIVER")
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_IP"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}
}

func refreshChannelTable() error {
	err := server.DB.DropTableIfExists(&models.Channel{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Channel{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func loadOneChannel() (models.Channel, error) {

	refreshChannelTable()

	channel := models.Channel{
		ID:          1,
		Channel:     "USSD",
		Description: "Test Channel",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := server.DB.Model(&models.Channel{}).Create(&channel).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return channel, nil
}

func refreshCustomerAndAccountTable() error {

	err := server.DB.DropTableIfExists(&models.Customer{}, &models.Accounts{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Customer{}, &models.Accounts{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func refreshRequestsAndTransactionsTable() error {

	err := server.DB.DropTableIfExists(&models.TransactionRequests{}, &models.Transactions{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.TransactionRequests{}, &models.Transactions{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func loadOneCustomer() (models.Customer, error) {

	err := refreshCustomerAndAccountTable()
	if err != nil {
		return models.Customer{}, err
	}
	customer := models.Customer{
		FirstName: "Cash",
		LastName:  "Account",
		Email:     "cash@gmail.com",
		AccountNo: "254712345678",
	}
	err = server.DB.Model(&models.Customer{}).Create(&customer).Error
	if err != nil {
		return models.Customer{}, err
	}
	account := models.Accounts{
		AccountNo:    "254712345678",
		Dormant:      "N",
		Approved:     true,
		ActualBal:    "1000.00",
		AvailableBal: "1000.00",
	}
	err = server.DB.Model(&models.Accounts{}).Create(&account).Error
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func loadCustomers() ([]models.Customer, error) {

	var err error

	if err != nil {
		return []models.Customer{}, err
	}
	var customers = []models.Customer{
		models.Customer{
			FirstName: "Cash",
			LastName:  "Account",
			Email:     "cash@gmail.com",
			AccountNo: "254712345678",
		},
		models.Customer{
			FirstName: "Simon",
			LastName:  "Maingi",
			Email:     "simon@gmail.com",
			AccountNo: "254708003472",
		},
	}
	var accounts = []models.Accounts{
		models.Accounts{
			AccountNo:    "254712345678",
			Dormant:      "N",
			Approved:     true,
			ActualBal:    "10000.00",
			AvailableBal: "10000.00",
		},
		models.Accounts{
			AccountNo:    "254708003472",
			Dormant:      "N",
			Approved:     true,
			ActualBal:    "100.00",
			AvailableBal: "100.00",
		},
	}

	for i := range customers {
		err = server.DB.Model(&models.Customer{}).Create(&customers[i]).Error
		if err != nil {
			log.Fatalf("cannot load customers table: %v", err)
		}

		err = server.DB.Model(&models.Accounts{}).Create(&accounts[i]).Error
		if err != nil {
			log.Fatalf("cannot load accounts table: %v", err)
		}
	}
	return customers, nil
}
