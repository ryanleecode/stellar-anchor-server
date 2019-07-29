module github.com/drdgvhbh/stellar-fi-anchor/ethereum

go 1.12

require (
	github.com/drdgvhbh/go-ethereum-hdwallet v0.0.0-20190717030924-e5b49683f92e
	github.com/drdgvhbh/stellar-fi-anchor v0.0.0-20190720024809-8522195fde21
	github.com/drdgvhbh/stellar-fi-anchor/middleware v0.0.0
	github.com/ethereum/go-ethereum v1.9.0
	github.com/go-errors/errors v1.0.1
	github.com/gorilla/handlers v1.4.1
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/schema v1.1.0
	github.com/jinzhu/gorm v1.9.10
	github.com/joho/godotenv v1.3.0
	github.com/pkg/errors v0.8.1
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stellar/go v0.0.0-20190719224126-f1f9d7901ae5
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/stretchr/testify v1.3.0
	github.com/thedevsaddam/govalidator v1.9.8
)

replace github.com/drdgvhbh/stellar-fi-anchor/middleware => ../middleware

replace github.com/drdgvhbh/stellar-fi-anchor/sdk => ../sdk
