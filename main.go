package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	router := gin.Default()
	router.GET("/customers", getCustomers)
	router.GET("/customers/:id", getCustomerById)
	router.POST("/customers", addCustomers)
	router.POST("/customers/login", login)
	router.POST("/payments", addPayment)
	router.GET("/payments", getPaymentHistory)

	router.Run(os.Getenv("DATABASE_HOST") + ":" + os.Getenv("DATABASE_PORT"))
}

type customer struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type payment struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	PaymentAmount string `json:"paymentAmount"`
}

type credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var customers = []customer{
	{ID: "1", Username: "customerone", Password: "pass1"},
	{ID: "2", Username: "customertwo", Password: "pass2"},
	{ID: "3", Username: "customerthree", Password: "pass3"},
}

var paymentHistory = []payment{}

func getCustomers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, customers)
}

func addCustomers(c *gin.Context) {
	var newCustomer customer

	if err := c.BindJSON(&newCustomer); err != nil {
		return
	}

	customers = append(customers, newCustomer)
	c.IndentedJSON(http.StatusCreated, newCustomer)
}

func getCustomerById(c *gin.Context) {
	id := c.Param("id")

	for _, customer := range customers {
		if customer.ID == id {
			c.IndentedJSON(http.StatusOK, customer)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "customer not found"})
}

func login(c *gin.Context) {
	var newCredential credential
	c.BindJSON(&newCredential)

	for _, customer := range customers {
		if newCredential.Username == customer.Username && newCredential.Password == customer.Password {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "successfully logged in"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "invalid username or password"})
}

func addPayment(c *gin.Context) {
	var newPayment payment

	if err := c.BindJSON(&newPayment); err != nil {
		return
	}

	for _, ph := range paymentHistory {
		if newPayment.Username == ph.Username {
			paymentHistory = append(paymentHistory, newPayment)
			c.IndentedJSON(http.StatusCreated, newPayment)
		}
	}
}

func getPaymentHistory(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, paymentHistory)
}
