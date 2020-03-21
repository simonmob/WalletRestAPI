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
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env,  %v", err)
	} else {
		fmt.Println("Loading the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_IP"), os.Getenv("DB_NAME"), os.Getenv("CASH_ACCOUNT"))

	sampledata.Load(server.DB)

	server.Run(":" + os.Getenv("API_PORT"))
}
