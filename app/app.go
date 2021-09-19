package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/wilker/banking/domain"
	"github.com/wilker/banking/service"
)

func sanityCheck() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Environment variable not defined")
	}
	envProps := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWD",
		"DB_ADDR",
		"DB_PORT",
		"DB_NAME",
		"AUTH_SERVER_ADDRESS",
		"AUTH_SERVER_PORT",
		"AUTH_PATH",
		"AUTH_SCHEME",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			log.Println(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}

func Start() {
	sanityCheck()

	router := mux.NewRouter()
	am := AuthMiddleware{domain.NewAuthRepository()}
	router.Use(am.authorizationHandler())

	dbClient := getDbClient()
	customerRepositoryDB := domain.NewCustomerRepositoryDb(dbClient)
	accountRepositoryDB := domain.NewAccountRepositoryDb(dbClient)

	// wiring
	// ch := CustomerHandlers{service: service.NewCustomerService(domain.NewCustomerRepositoryStub())}
	ch := CustomerHandlers{service: service.NewCustomerService(customerRepositoryDB)}
	ah := AccountHandler{service: service.NewAccountService(accountRepositoryDB)}

	// define routes
	router.HandleFunc("/customers", ch.getAllCustomers).
		Methods(http.MethodGet).
		Name("GetAllCustomers")
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomer).
		Methods(http.MethodGet).
		Name("GetCustomer")
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account", ah.NewAccount).
		Methods(http.MethodPost).
		Name("NewAccount")
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account/{account_id:[0-9]+}", ah.MakeTransaction).
		Methods(http.MethodPost).
		Name("NewTransaction")

	// starting server
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}

func getDbClient() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbAddr, dbPort, dbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}
