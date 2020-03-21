package middlewares

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"tospay.com/WalletRestAPI/api/authorization"
	"tospay.com/WalletRestAPI/api/responses"
)

//SetMiddlewareJSON adds Content type on header
func SetMiddlewareJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

//SetMiddlewareAuthentication adds custom jwt authentication middlewares
func SetMiddlewareAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := authorization.TokenValid(c)
		if err != nil {
			responses.ERROR(c, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		c.Next()
	}
}

//SetMiddlewareLogger set logger to format logs
func SetMiddlewareLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()

		// Logging to a file.
		f, _ := os.Create("walletapi.log")
		gin.DefaultWriter = io.MultiWriter(f)

		// Use the following code if you need to write the logs to file and console at the same time.
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		c.Next()
	}

}
