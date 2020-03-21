package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //Postgres database DRIVER
)

//Server struct holds DB and Gin Router details
type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
	Configs
}

//Configs holds other configs in .ENV file
type Configs struct {
	CashAccount string //used as settlement account for CashDeposit and cashwithdrawals
}

//Initialize initializes postgres database and gin Router
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName, CashAccount string) {

	var err error

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}

	//Create all the tables
	//server.DB.Debug().AutoMigrate(&models.Channel{}, &models.Accounts{}, &models.Customer{}, &models.TransactionRequests{}, &models.Transactions{}) //database migration

	server.Configs.CashAccount = CashAccount

	server.Router = gin.New()
	server.initializeRoutes() //initializes all routes created in routes.go
}

//Run starta the api Listening on the given address. default is 8080
func (server *Server) Run(addr string) {
	fmt.Println("Listening on port " + addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
