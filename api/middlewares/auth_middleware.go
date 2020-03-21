package middlewares

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
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
	//return func(c *gin.Context) {
	// Disable Console Color, you don't need console color when writing the logs to file.
	// gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("walletapi")
	gin.DefaultWriter = io.MultiWriter(f)

	// // Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	// c.Next()

	// write file
	fileName := "walletapi"
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	// instantiation
	logger := logrus.New()
	// Set output
	logger.Out = src
	// Set log level
	logger.SetLevel(logrus.DebugLevel)
	// Set rotatelogs
	logWriter, _ := rotatelogs.New(
		// Split file name
		fileName+".%Y%m%d.log",
		// Generate soft chain, point to the latest log file
		rotatelogs.WithLinkName(fileName),
		// Set maximum save time (7 days)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// Set log cutting interval (1 day)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// Add Hook
	logger.AddHook(lfHook)
	return func(c *gin.Context) {
		// start time
		startTime := time.Now()
		// Processing request
		c.Next()
		// End time
		endTime := time.Now()
		// execution time
		latencyTime := endTime.Sub(startTime)
		// Request mode
		reqMethod := c.Request.Method
		// Request routing
		reqURI := c.Request.RequestURI
		// Status code
		statusCode := c.Writer.Status()
		// Request IP
		clientIP := c.ClientIP()
		// Log format
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqURI,
		}).Info()
	}
}
