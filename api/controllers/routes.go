package controllers

import (
	"github.com/gin-gonic/gin"
	"tospay.com/WalletRestAPI/api/middlewares"
)

//initializeRoutes initializes routes for incoming Requests
func (s *Server) initializeRoutes() {

	s.Router.Use(middlewares.SetMiddlewareJSON())   //adds json content type on request
	s.Router.Use(middlewares.SetMiddlewareLogger()) //Add logger middleware
	s.Router.Use(gin.Recovery())                    //recovers from crashes

	// Home Route
	s.Router.GET("/", s.Home)

	// Channel authorization routes
	s.Router.POST("/channelAuth", s.ChannelAuth)

	//channel routes
	//channelsGroup := s.Router.Group("/channels")
	s.Router.POST("/channels", middlewares.SetMiddlewareLogger(), s.CreateChannel)
	s.Router.GET("/GetChannels", s.GetChannels)

	//Customer routes
	customersGroup := s.Router.Group("/customers")
	customersGroup.Use(middlewares.SetMiddlewareAuthentication()) //must be authenticated.
	{
		customersGroup.POST("/create", s.CreateCustomer)
		customersGroup.GET("/get", s.GetCustomers)
		customersGroup.GET("/get/:accountno", s.GetCustomer)
		customersGroup.PUT("/update/:id", s.UpdateCustomer)

	}
	//Transactions route
	s.Router.POST("/transactions", middlewares.SetMiddlewareAuthentication(), s.IncomingRequest)

}
