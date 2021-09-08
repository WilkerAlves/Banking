package domain

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/wilker/banking/errs"
	"github.com/wilker/banking/logger"
	"time"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {
	var err error
	customers := make([]Customer, 0)

	if status == "" {
		findAllSQL := `select customer_id, name, city, zipcode, date_of_birth, status from banking.customers`
		err = d.client.Select(&customers, findAllSQL)
	} else {
		findAllSQL := `select customer_id, name, city, zipcode, date_of_birth, status from banking.customers where status = ?`
		err = d.client.Select(&customers, findAllSQL, status)
	}

	if err != nil {
		logger.Error("Error while querying customers table " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return customers, nil
}

func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	var c Customer
	customerSql := `select customer_id, name, city, zipcode, date_of_birth, status 
					from banking.customers 
					where customer_id = ?`
	err := d.client.Get(&c, customerSql, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("customer not found")
		} else {
			logger.Error("Error while scanning customer " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}
	return &c, nil
}

func NewCustomerRepositoryDb() CustomerRepositoryDb {
	client, err := sqlx.Open("mysql", "root:4874021993@tcp(localhost:3306)/banking")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return CustomerRepositoryDb{client: client}
}
