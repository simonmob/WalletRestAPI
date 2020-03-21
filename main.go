package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"tospay.com/WalletRestAPI/api/controllers"
	"tospay.com/WalletRestAPI/api/sampledata"
)

var server = controllers.Server{}

func main() {

	var err error
	err = godotenv.Load() //logs the .env for Configs
	if err != nil {
		log.Fatalf("Error getting env,  %v", err)
	} else {
		fmt.Println("Loading the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_IP"), os.Getenv("DB_NAME"), os.Getenv("CASH_ACCOUNT"))

	sampledata.Load(server.DB) //AutoMigrate all the models and insert test ssample data.(USSD chaneel and the cash ACCOUNT)

	server.Run(":" + os.Getenv("API_PORT"))
}
