/*
 * Copyright (c) 2021.
 * Marc Concepcion
 * marcanthonyconcepcion@gmail.com
 */

package MarcGoRESTAPIDemo

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

var settings = readConfiguration("resources/MarcGoRESTAPIDemo.yaml")

type Records struct {
	database *sql.DB
}

type Subscriber struct {
	index          uint8
	emailAddress   string
	firstName      string
	lastName       string
	activationFlag bool
}

func makeDatabaseRecords() Records {
	database, dbInstanceFail := sql.Open("mysql", settings.Database.User+":"+settings.Database.Password+
		"@tcp("+settings.Database.Host+":"+strconv.Itoa(int(settings.Database.Port))+")/"+settings.Database.DBName)
	if dbInstanceFail != nil {
		panic(dbInstanceFail.Error())
	}
	return Records{database}
}

func (records Records) create(subscriber Subscriber) (sql.Result, error) {
	result, fault := records.database.Exec(
		"insert into `subscribers` (`email_address`, `last_name`, `first_name`) values (?, ?, ?)",
		subscriber.emailAddress, subscriber.lastName, subscriber.firstName)
	return result, fault
}

func (records Records) retrieve(index uint8) (*Subscriber, error) {
	var subscriber Subscriber
	record := records.database.QueryRow("select * from `subscribers` where `index`=?", index)
	if recordModelError := record.Scan(&subscriber.index, &subscriber.emailAddress, &subscriber.lastName, &subscriber.firstName,
		&subscriber.activationFlag); recordModelError != nil {
		return &subscriber, recordModelError
	}
	return &subscriber, nil
}

func (records Records) update(index uint8, updateValues map[string]string) (sql.Result, error) {
	var parametersToUpdate []string
	emailAddress, emailAddressKey := updateValues["email_address"]
	if true == emailAddressKey {
		parametersToUpdate = append(parametersToUpdate, "`email_address` = "+"\""+emailAddress+"\"")
	}
	lastName, lastNameKey := updateValues["last_name"]
	if true == lastNameKey {
		parametersToUpdate = append(parametersToUpdate, "`last_name` = "+"\""+lastName+"\"")
	}
	firstName, firstNameKey := updateValues["first_name"]
	if true == firstNameKey {
		parametersToUpdate = append(parametersToUpdate, "`first_name` = "+"\""+firstName+"\"")
	}
	activationFlag, activationFlagKey := updateValues["activation_flag"]
	if true == activationFlagKey {
		parametersToUpdate = append(parametersToUpdate, "`activation_flag` = "+activationFlag)
	}
	result, updateFail := records.database.Exec("update `subscribers` set "+
		strings.Join(parametersToUpdate, ",")+" where `index`=?", index)
	return result, updateFail
}

func (records Records) delete(index uint8) (sql.Result, error) {
	result, deleteError := records.database.Exec("delete from `subscribers` where `index`=?", index)
	return result, deleteError
}

func (records Records) list() ([]Subscriber, error) {
	rows, dbQueryError := records.database.Query("select * from `subscribers`")
	if dbQueryError != nil {
		return nil, dbQueryError
	}
	subscribers := make([]Subscriber, 0)
	for rows.Next() {
		var subscriber Subscriber
		if recordModelError := rows.Scan(&subscriber.index, &subscriber.emailAddress, &subscriber.lastName, &subscriber.firstName,
			&subscriber.activationFlag); recordModelError != nil {
			return subscribers, recordModelError
		}
		subscribers = append(subscribers, subscriber)
	}
	rowsCloseError := rows.Close()
	if rowsCloseError != nil {
		return subscribers, rowsCloseError
	}
	return subscribers, nil
}
